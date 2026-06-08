'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { api } from '@/lib/api'
import type { Item, ItemStatus, ItemPriority } from '@/types'
import { STATUS_LABELS, PRIORITY_LABELS } from '@/types'
import { Plus, Search, Filter, ChevronRight, Circle } from 'lucide-react'

export const runtime = 'edge'

const STATUS_COLORS: Record<string, string> = {
  draft:       'bg-[#4B5563]/20 text-[#9CA3AF] border-[#4B5563]/30',
  pending:     'bg-[#D97706]/20 text-[#FCD34D] border-[#D97706]/30',
  in_progress: 'bg-[#2563EB]/20 text-[#60A5FA] border-[#2563EB]/30',
  waiting:     'bg-[#7C3AED]/20 text-[#A78BFA] border-[#7C3AED]/30',
  done:        'bg-[#059669]/20 text-[#34D399] border-[#059669]/30',
  rejected:    'bg-[#DC2626]/20 text-[#F87171] border-[#DC2626]/30',
}

const PRIORITY_COLORS: Record<string, string> = {
  low:    'text-[#6B7280]',
  medium: 'text-[#60A5FA]',
  high:   'text-[#FCD34D]',
  urgent: 'text-[#F87171]',
}

const ALL_STATUSES = Object.entries(STATUS_LABELS) as [ItemStatus, string][]

export default function TasksPage() {
  const [items, setItems] = useState<Item[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [filterStatus, setFilterStatus] = useState<ItemStatus | ''>('')

  useEffect(() => {
    api.items.list()
      .then(r => setItems(r.items ?? []))
      .catch(() => setItems([]))
      .finally(() => setLoading(false))
  }, [])

  const rootItems = items.filter(i => !i.parent_id)

  const filtered = rootItems.filter(i => {
    if (filterStatus && i.status !== filterStatus) return false
    if (search && !i.title.toLowerCase().includes(search.toLowerCase())) return false
    return true
  })

  const getChildren = (parentId: string) => items.filter(i => i.parent_id === parentId)

  return (
    <div className="space-y-5 animate-fade-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="font-display text-xl font-bold text-[#E8EEF4]">Tasks & Requisições</h1>
          <p className="text-[#4A5C6A] text-sm mt-0.5">{filtered.length} item{filtered.length !== 1 ? 's' : ''}</p>
        </div>
        <Link
          href="/tasks/new"
          className="flex items-center gap-2 bg-brand-cyan text-[#080D12] font-display font-semibold text-sm px-4 py-2 rounded-lg hover:bg-white transition-colors shadow-accent"
        >
          <Plus size={14} />
          Novo Item
        </Link>
      </div>

      {/* Filters */}
      <div className="flex gap-3 flex-wrap">
        <div className="relative flex-1 min-w-[200px]">
          <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-[#4A5C6A]" />
          <input
            value={search}
            onChange={e => setSearch(e.target.value)}
            placeholder="Buscar items..."
            className="w-full bg-[#0D1620] border border-white/[0.07] rounded-lg pl-9 pr-4 py-2 text-sm text-[#E8EEF4] placeholder-[#4A5C6A] outline-none focus:border-brand-cyan/40 transition-colors"
          />
        </div>
        <div className="relative">
          <Filter size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-[#4A5C6A] pointer-events-none" />
          <select
            value={filterStatus}
            onChange={e => setFilterStatus(e.target.value as ItemStatus | '')}
            className="bg-[#0D1620] border border-white/[0.07] rounded-lg pl-9 pr-8 py-2 text-sm text-[#E8EEF4] outline-none focus:border-brand-cyan/40 transition-colors appearance-none cursor-pointer"
          >
            <option value="">Todos os status</option>
            {ALL_STATUSES.map(([k, v]) => <option key={k} value={k}>{v}</option>)}
          </select>
        </div>
      </div>

      {/* Item list */}
      {loading ? (
        <div className="flex justify-center py-16">
          <div className="w-6 h-6 border-2 border-brand-cyan/40 border-t-brand-cyan rounded-full animate-spin" />
        </div>
      ) : filtered.length === 0 ? (
        <div className="card p-12 text-center">
          <p className="text-[#4A5C6A]">Nenhum item encontrado.</p>
        </div>
      ) : (
        <div className="space-y-2">
          {filtered.map(item => {
            const children = getChildren(item.id)
            return (
              <div key={item.id} className="group">
                <Link
                  href={`/tasks/${item.id}`}
                  className="card flex items-center gap-4 px-5 py-4 hover:border-white/[0.12] hover:bg-[#0D1620]/80 transition-all"
                >
                  {/* Priority dot */}
                  <Circle
                    size={8}
                    className={`shrink-0 ${PRIORITY_COLORS[item.priority]} fill-current`}
                  />

                  {/* Title + meta */}
                  <div className="flex-1 min-w-0">
                    <p className="text-sm text-[#E8EEF4] font-medium truncate group-hover:text-brand-cyan transition-colors">
                      {item.title}
                    </p>
                    <div className="flex items-center gap-3 mt-1">
                      {item.deadline && (
                        <span className="text-[11px] text-[#4A5C6A] font-mono">
                          prazo {new Date(item.deadline).toLocaleDateString('pt-BR')}
                        </span>
                      )}
                      {children.length > 0 && (
                        <span className="text-[11px] text-[#4A5C6A] font-mono">
                          {children.length} sub-item{children.length !== 1 ? 's' : ''}
                        </span>
                      )}
                      <span className={`text-[11px] font-mono ${PRIORITY_COLORS[item.priority]}`}>
                        {PRIORITY_LABELS[item.priority]}
                      </span>
                    </div>
                  </div>

                  {/* Status pill */}
                  <span className={`status-pill border ${STATUS_COLORS[item.status] ?? ''} shrink-0`}>
                    {STATUS_LABELS[item.status]}
                  </span>

                  <ChevronRight size={14} className="text-[#4A5C6A] group-hover:text-brand-cyan shrink-0 transition-colors" />
                </Link>

                {/* Children preview */}
                {children.length > 0 && (
                  <div className="ml-6 pl-4 border-l border-white/[0.04] space-y-0.5 mt-0.5">
                    {children.slice(0, 3).map(child => (
                      <Link
                        key={child.id}
                        href={`/tasks/${child.id}`}
                        className="flex items-center gap-3 px-4 py-2 rounded-lg hover:bg-white/[0.02] transition-colors group/child"
                      >
                        <Circle size={6} className={`shrink-0 ${PRIORITY_COLORS[child.priority]} fill-current opacity-60`} />
                        <span className="text-xs text-[#7D95A8] truncate group-hover/child:text-[#E8EEF4] transition-colors flex-1">{child.title}</span>
                        <span className={`status-pill border text-[10px] ${STATUS_COLORS[child.status] ?? ''}`}>
                          {STATUS_LABELS[child.status]}
                        </span>
                      </Link>
                    ))}
                    {children.length > 3 && (
                      <p className="text-[11px] text-[#4A5C6A] pl-4 py-1 font-mono">+{children.length - 3} mais…</p>
                    )}
                  </div>
                )}
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
