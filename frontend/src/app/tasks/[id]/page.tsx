'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import { api } from '@/lib/api'
import type { Item, Comment } from '@/types'
import { STATUS_LABELS, PRIORITY_LABELS } from '@/types'
import { ArrowLeft, MessageSquare, AlertCircle, ChevronRight, Send } from 'lucide-react'

const STATUS_COLORS: Record<string, string> = {
  draft:       'bg-[#4B5563]/20 text-[#9CA3AF] border-[#4B5563]/30',
  pending:     'bg-[#D97706]/20 text-[#FCD34D] border-[#D97706]/30',
  in_progress: 'bg-[#2563EB]/20 text-[#60A5FA] border-[#2563EB]/30',
  waiting:     'bg-[#7C3AED]/20 text-[#A78BFA] border-[#7C3AED]/30',
  done:        'bg-[#059669]/20 text-[#34D399] border-[#059669]/30',
  rejected:    'bg-[#DC2626]/20 text-[#F87171] border-[#DC2626]/30',
}

const NEXT_STATUS: Record<string, string> = {
  draft:       'pending',
  pending:     'in_progress',
  in_progress: 'done',
}

const NEXT_LABEL: Record<string, string> = {
  draft:       'Submeter',
  pending:     'Iniciar',
  in_progress: 'Concluir',
}

export default function TaskDetailPage() {
  const { id } = useParams<{ id: string }>()
  const router = useRouter()
  const [item, setItem] = useState<Item | null>(null)
  const [loading, setLoading] = useState(true)
  const [comment, setComment] = useState('')
  const [posting, setPosting] = useState(false)

  async function load() {
    try {
      const data = await api.items.get(id)
      setItem(data)
    } catch {
      router.push('/tasks')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { load() }, [id])

  async function advanceStatus() {
    if (!item || !NEXT_STATUS[item.status]) return
    const updated = await api.items.update(id, { status: NEXT_STATUS[item.status] as any })
    setItem(updated)
  }

  async function submitComment() {
    if (!comment.trim() || posting) return
    setPosting(true)
    try {
      await api.items.comments.create(id, comment.trim())
      setComment('')
      load()
    } finally {
      setPosting(false)
    }
  }

  if (loading) return (
    <div className="flex justify-center py-16">
      <div className="w-6 h-6 border-2 border-brand-cyan/40 border-t-brand-cyan rounded-full animate-spin" />
    </div>
  )

  if (!item) return null

  return (
    <div className="max-w-3xl space-y-5 animate-fade-in">
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-sm text-[#4A5C6A]">
        <Link href="/tasks" className="hover:text-brand-cyan transition-colors flex items-center gap-1">
          <ArrowLeft size={13} /> Tasks
        </Link>
        <ChevronRight size={12} />
        <span className="text-[#7D95A8] truncate max-w-xs">{item.title}</span>
      </div>

      {/* Header card */}
      <div className="card card-accent p-6">
        <div className="flex items-start justify-between gap-4 mb-4">
          <h1 className="font-display text-xl font-bold text-[#E8EEF4] flex-1">{item.title}</h1>
          <span className={`status-pill border ${STATUS_COLORS[item.status] ?? ''} shrink-0 mt-1`}>
            {STATUS_LABELS[item.status]}
          </span>
        </div>

        {item.description && (
          <p className="text-[#7D95A8] text-sm leading-relaxed mb-4">{item.description}</p>
        )}

        {/* Meta grid */}
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 text-xs">
          {[
            { label: 'Prioridade', value: PRIORITY_LABELS[item.priority] },
            { label: 'Visibilidade', value: item.visibility === 'private' ? 'Privado' : item.visibility === 'team' ? 'Equipe' : 'Público' },
            { label: 'Criado em', value: new Date(item.created_at).toLocaleDateString('pt-BR') },
            item.deadline ? { label: 'Prazo', value: new Date(item.deadline).toLocaleDateString('pt-BR') } : null,
          ].filter(Boolean).map(({ label, value }: any) => (
            <div key={label} className="bg-white/[0.03] rounded-lg px-3 py-2.5 border border-white/[0.05]">
              <p className="text-[#4A5C6A] font-mono uppercase tracking-widest mb-1" style={{ fontSize: 9 }}>{label}</p>
              <p className="text-[#E8EEF4] font-medium">{value}</p>
            </div>
          ))}
        </div>

        {/* Actions */}
        {NEXT_STATUS[item.status] && (
          <div className="flex gap-3 mt-5 pt-4 border-t border-white/[0.06]">
            <button
              onClick={advanceStatus}
              className="flex items-center gap-2 bg-brand-cyan text-[#080D12] font-display font-semibold text-sm px-4 py-2 rounded-lg hover:bg-white transition-colors"
            >
              {NEXT_LABEL[item.status]}
            </button>
          </div>
        )}
      </div>

      {/* Sub-items */}
      {item.children && item.children.length > 0 && (
        <div className="card p-0 overflow-hidden">
          <div className="px-5 py-3.5 border-b border-white/[0.06]">
            <h2 className="font-display font-semibold text-sm text-[#E8EEF4]">
              Sub-items <span className="text-[#4A5C6A] font-normal">({item.children.length})</span>
            </h2>
          </div>
          <div className="divide-y divide-white/[0.04]">
            {item.children.map(child => (
              <Link key={child.id} href={`/tasks/${child.id}`}
                className="flex items-center gap-3 px-5 py-3 hover:bg-white/[0.02] transition-colors group"
              >
                <span className="text-sm text-[#E8EEF4] flex-1 truncate group-hover:text-brand-cyan transition-colors">{child.title}</span>
                <span className={`status-pill border ${STATUS_COLORS[child.status] ?? ''}`}>{STATUS_LABELS[child.status]}</span>
              </Link>
            ))}
          </div>
        </div>
      )}

      {/* Comments */}
      <div className="card p-0 overflow-hidden">
        <div className="flex items-center justify-between px-5 py-3.5 border-b border-white/[0.06]">
          <h2 className="font-display font-semibold text-sm text-[#E8EEF4] flex items-center gap-2">
            <MessageSquare size={14} className="text-[#4A5C6A]" />
            Comentários <span className="text-[#4A5C6A] font-normal">({item.comments?.length ?? 0})</span>
          </h2>
        </div>

        {/* Comment list */}
        <div className="divide-y divide-white/[0.04]">
          {(item.comments ?? []).length === 0 ? (
            <p className="px-5 py-4 text-[#4A5C6A] text-sm">Nenhum comentário ainda.</p>
          ) : (
            item.comments!.map(c => (
              <div key={c.id} className="px-5 py-4">
                <div className="flex items-center gap-2 mb-1.5">
                  <div className="w-5 h-5 rounded-full bg-brand-cyan/20 border border-brand-cyan/30 flex items-center justify-center text-brand-cyan text-[9px] font-display font-bold">
                    {c.author?.name?.[0] ?? '?'}
                  </div>
                  <span className="text-xs text-[#7D95A8] font-medium">{c.author?.name ?? 'Usuário'}</span>
                  <span className="text-[10px] text-[#4A5C6A] font-mono ml-auto">{new Date(c.created_at).toLocaleString('pt-BR')}</span>
                </div>
                <p className="text-sm text-[#E8EEF4] leading-relaxed pl-7">{c.content}</p>
              </div>
            ))
          )}
        </div>

        {/* Comment input */}
        <div className="px-5 py-4 border-t border-white/[0.06] flex gap-3">
          <input
            value={comment}
            onChange={e => setComment(e.target.value)}
            onKeyDown={e => e.key === 'Enter' && !e.shiftKey && submitComment()}
            placeholder="Adicionar comentário..."
            className="flex-1 bg-[#0D1620] border border-white/[0.07] rounded-lg px-4 py-2 text-sm text-[#E8EEF4] placeholder-[#4A5C6A] outline-none focus:border-brand-cyan/40 transition-colors"
          />
          <button
            onClick={submitComment}
            disabled={!comment.trim() || posting}
            className="p-2 rounded-lg bg-brand-cyan/10 border border-brand-cyan/30 text-brand-cyan hover:bg-brand-cyan hover:text-[#080D12] transition-colors disabled:opacity-40"
          >
            <Send size={14} />
          </button>
        </div>
      </div>
    </div>
  )
}
