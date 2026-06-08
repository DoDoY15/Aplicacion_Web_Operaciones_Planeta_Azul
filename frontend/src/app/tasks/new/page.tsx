'use client'

import { useState, FormEvent } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { api } from '@/lib/api'
import type { ItemVisibility, ItemPriority } from '@/types'
import { ArrowLeft } from 'lucide-react'

export default function NewTaskPage() {
  const router = useRouter()
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [visibility, setVisibility] = useState<ItemVisibility>('team')
  const [priority, setPriority] = useState<ItemPriority>('medium')
  const [deadline, setDeadline] = useState('')
  const [requiresApproval, setRequiresApproval] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setLoading(true)
    setError('')
    try {
      const item = await api.items.create({
        title, description, visibility, priority,
        requires_approval: requiresApproval,
        deadline: deadline ? new Date(deadline).toISOString() : undefined,
      })
      router.push(`/tasks/${item.id}`)
    } catch (err: any) {
      setError(err.message ?? 'Error al crear el ítem')
    } finally {
      setLoading(false)
    }
  }

  const inputClass = "w-full bg-[#0D1620] border border-white/[0.07] rounded-lg px-4 py-2.5 text-sm text-[#E8EEF4] placeholder-[#4A5C6A] outline-none focus:border-brand-cyan/40 transition-colors"
  const labelClass = "block text-[10px] font-mono text-[#7D95A8] uppercase tracking-widest mb-1.5"

  return (
    <div className="max-w-2xl space-y-5 animate-fade-in">
      <div className="flex items-center gap-3">
        <Link href="/tasks" className="text-[#4A5C6A] hover:text-brand-cyan transition-colors">
          <ArrowLeft size={16} />
        </Link>
        <h1 className="font-display text-xl font-bold text-[#E8EEF4]">Nuevo Ítem</h1>
      </div>

      <form onSubmit={handleSubmit} className="card card-accent p-6 space-y-5">
        <div>
          <label className={labelClass}>Título *</label>
          <input value={title} onChange={e => setTitle(e.target.value)}
            placeholder="Describe el ítem en pocas palabras..."
            required minLength={3} maxLength={255} className={inputClass} />
        </div>

        <div>
          <label className={labelClass}>Descripción</label>
          <textarea value={description} onChange={e => setDescription(e.target.value)}
            placeholder="Contexto, detalles, criterios de finalización..."
            rows={4}
            className={`${inputClass} resize-none`} />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className={labelClass}>Visibilidad</label>
            <select value={visibility} onChange={e => setVisibility(e.target.value as ItemVisibility)}
              className={`${inputClass} cursor-pointer`}>
              <option value="private">Privado</option>
              <option value="team">Equipo</option>
              <option value="public">Público</option>
            </select>
          </div>
          <div>
            <label className={labelClass}>Prioridad</label>
            <select value={priority} onChange={e => setPriority(e.target.value as ItemPriority)}
              className={`${inputClass} cursor-pointer`}>
              <option value="low">Baja</option>
              <option value="medium">Media</option>
              <option value="high">Alta</option>
              <option value="urgent">Urgente</option>
            </select>
          </div>
        </div>

        <div>
          <label className={labelClass}>Fecha límite (opcional)</label>
          <input type="datetime-local" value={deadline} onChange={e => setDeadline(e.target.value)}
            className={inputClass} />
        </div>

        <label className="flex items-center gap-3 cursor-pointer group">
          <div
            onClick={() => setRequiresApproval(v => !v)}
            className={`w-9 h-5 rounded-full border transition-colors flex items-center px-0.5 ${requiresApproval ? 'bg-brand-cyan/20 border-brand-cyan/50' : 'bg-white/[0.04] border-white/10'}`}
          >
            <div className={`w-4 h-4 rounded-full transition-transform ${requiresApproval ? 'bg-brand-cyan translate-x-4' : 'bg-[#4A5C6A] translate-x-0'}`} />
          </div>
          <span className="text-sm text-[#7D95A8]">Requiere aprobación antes de iniciar</span>
        </label>

        {error && (
          <div className="bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-2.5 text-red-400 text-sm">
            {error}
          </div>
        )}

        <div className="flex gap-3 pt-2">
          <button type="submit" disabled={loading}
            className="flex-1 bg-brand-cyan text-[#080D12] font-display font-semibold text-sm py-2.5 rounded-lg hover:bg-white transition-colors disabled:opacity-50">
            {loading ? 'Creando...' : 'Crear Ítem'}
          </button>
          <Link href="/tasks"
            className="px-5 py-2.5 rounded-lg border border-white/[0.07] text-[#7D95A8] text-sm hover:border-white/20 hover:text-[#E8EEF4] transition-colors text-center">
            Cancelar
          </Link>
        </div>
      </form>
    </div>
  )
}
