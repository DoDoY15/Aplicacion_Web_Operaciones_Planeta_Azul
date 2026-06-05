'use client'

import { useEffect, useState } from 'react'
import { api } from '@/lib/api'
import type { Item } from '@/types'
import { STATUS_LABELS, PRIORITY_LABELS } from '@/types'
import { ListTodo, CheckCircle2, Clock, AlertTriangle, TrendingUp } from 'lucide-react'
import Link from 'next/link'

const STATUS_COLORS: Record<string, string> = {
  draft:       'bg-[#4B5563]/20 text-[#9CA3AF]',
  pending:     'bg-[#D97706]/20 text-[#FCD34D]',
  in_progress: 'bg-[#2563EB]/20 text-[#60A5FA]',
  waiting:     'bg-[#7C3AED]/20 text-[#A78BFA]',
  done:        'bg-[#059669]/20 text-[#34D399]',
  rejected:    'bg-[#DC2626]/20 text-[#F87171]',
}

const PRIORITY_DOTS: Record<string, string> = {
  low:    'bg-[#6B7280]',
  medium: 'bg-[#2563EB]',
  high:   'bg-[#D97706]',
  urgent: 'bg-[#DC2626] animate-pulse-slow',
}

export default function DashboardPage() {
  const [items, setItems] = useState<Item[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.items.list()
      .then(r => setItems(r.items ?? []))
      .catch(() => setItems([]))
      .finally(() => setLoading(false))
  }, [])

  const stats = {
    total:       items.length,
    in_progress: items.filter(i => i.status === 'in_progress').length,
    pending:     items.filter(i => i.status === 'pending').length,
    done:        items.filter(i => i.status === 'done').length,
    urgent:      items.filter(i => i.priority === 'urgent' || i.priority === 'high').length,
  }

  const recent = [...items]
    .sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
    .slice(0, 8)

  return (
    <div className="space-y-6 animate-fade-in">
      {/* Stats row */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        {[
          { label: 'Total de Items', value: stats.total, icon: ListTodo, color: 'text-brand-cyan' },
          { label: 'Em Andamento',   value: stats.in_progress, icon: TrendingUp, color: 'text-[#60A5FA]' },
          { label: 'Pendentes',      value: stats.pending, icon: Clock, color: 'text-[#FCD34D]' },
          { label: 'Alta Prioridade',value: stats.urgent, icon: AlertTriangle, color: 'text-[#F87171]' },
        ].map(({ label, value, icon: Icon, color }) => (
          <div key={label} className="card card-accent p-5">
            <div className="flex items-start justify-between">
              <div>
                <p className="text-[#4A5C6A] text-xs font-mono uppercase tracking-widest mb-2">{label}</p>
                <p className={`font-display text-3xl font-bold ${color}`}>
                  {loading ? '—' : value}
                </p>
              </div>
              <Icon size={18} className={`${color} opacity-60 mt-0.5`} />
            </div>
          </div>
        ))}
      </div>

      {/* Status breakdown */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        {/* Recent items */}
        <div className="card col-span-2 p-0 overflow-hidden">
          <div className="flex items-center justify-between px-5 py-4 border-b border-white/[0.06]">
            <h2 className="font-display font-semibold text-sm text-[#E8EEF4]">Atividade Recente</h2>
            <Link href="/tasks" className="text-brand-cyan text-xs font-mono hover:underline">Ver tudo →</Link>
          </div>
          {loading ? (
            <div className="p-8 flex justify-center">
              <div className="w-5 h-5 border-2 border-brand-cyan/40 border-t-brand-cyan rounded-full animate-spin" />
            </div>
          ) : recent.length === 0 ? (
            <div className="p-8 text-center text-[#4A5C6A] text-sm">Nenhum item ainda.</div>
          ) : (
            <div className="divide-y divide-white/[0.04]">
              {recent.map(item => (
                <Link
                  key={item.id}
                  href={`/tasks/${item.id}`}
                  className="flex items-center gap-4 px-5 py-3.5 hover:bg-white/[0.02] transition-colors group"
                >
                  <div className={`w-1.5 h-1.5 rounded-full shrink-0 ${PRIORITY_DOTS[item.priority]}`} />
                  <div className="flex-1 min-w-0">
                    <p className="text-sm text-[#E8EEF4] truncate group-hover:text-brand-cyan transition-colors">{item.title}</p>
                    {item.deadline && (
                      <p className="text-[11px] text-[#4A5C6A] font-mono mt-0.5">
                        prazo: {new Date(item.deadline).toLocaleDateString('pt-BR')}
                      </p>
                    )}
                  </div>
                  <span className={`status-pill ${STATUS_COLORS[item.status] ?? 'bg-white/10 text-[#7D95A8]'} shrink-0`}>
                    {STATUS_LABELS[item.status]}
                  </span>
                </Link>
              ))}
            </div>
          )}
        </div>

        {/* Status summary */}
        <div className="card p-5">
          <h2 className="font-display font-semibold text-sm text-[#E8EEF4] mb-4">Por Status</h2>
          <div className="space-y-3">
            {Object.entries(STATUS_LABELS).map(([status, label]) => {
              const count = items.filter(i => i.status === status).length
              const pct = items.length > 0 ? (count / items.length) * 100 : 0
              return (
                <div key={status}>
                  <div className="flex justify-between text-xs mb-1">
                    <span className="text-[#7D95A8]">{label}</span>
                    <span className="font-mono text-[#4A5C6A]">{count}</span>
                  </div>
                  <div className="h-1 bg-white/[0.06] rounded-full overflow-hidden">
                    <div
                      className="h-full rounded-full transition-all duration-500"
                      style={{
                        width: `${pct}%`,
                        background: status === 'done' ? '#059669' : status === 'in_progress' ? '#2563EB' : status === 'pending' ? '#D97706' : '#4B5563'
                      }}
                    />
                  </div>
                </div>
              )
            })}
          </div>
        </div>
      </div>
    </div>
  )
}
