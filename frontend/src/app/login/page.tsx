'use client'

import { useState, FormEvent } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/hooks/useAuth'
import { Logo } from '@/components/ui/Logo'

export default function LoginPage() {
  const { login } = useAuth()
  const router = useRouter()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await login(email, password)
      router.push('/dashboard')
    } catch {
      setError('Credenciales inválidas. Verifica el correo electrónico y la contraseña.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-[#080D12] bg-grid flex items-center justify-center p-4 relative overflow-hidden">
      {/* Ambient glow */}
      <div className="absolute inset-0 bg-glow-cyan pointer-events-none" />
      <div className="absolute top-1/3 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] rounded-full bg-brand-cyan/5 blur-3xl pointer-events-none" />

      <div className="relative w-full max-w-sm animate-fade-in">
        {/* Logo */}
        <div className="flex flex-col items-center mb-10">
          <Logo size={72} className="mb-4 rounded-full shadow-accent" />
          <h1 className="font-display text-2xl font-bold text-[#E8EEF4] tracking-tight">
            Planeta Azul
          </h1>
          <p className="text-[#4A5C6A] text-sm mt-1 font-mono">SISTEMA DE GESTIÓN INDUSTRIAL</p>
        </div>

        {/* Card */}
        <div className="card card-accent p-7 shadow-accent">
          <h2 className="font-display text-lg font-semibold text-[#E8EEF4] mb-6">
            Iniciar sesión en el sistema
          </h2>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-xs font-mono text-[#7D95A8] uppercase tracking-widest mb-1.5">
                Correo electrónico
              </label>
              <input
                type="email"
                value={email}
                onChange={e => setEmail(e.target.value)}
                placeholder="tu@email.com"
                required
                className="w-full bg-[#0D1620] border border-white/10 rounded-lg px-4 py-2.5 text-sm text-[#E8EEF4] placeholder-[#4A5C6A] outline-none focus:border-brand-cyan/50 focus:shadow-[0_0_0_3px_rgba(0,212,255,0.08)] transition-all"
              />
            </div>

            <div>
              <label className="block text-xs font-mono text-[#7D95A8] uppercase tracking-widest mb-1.5">
                Contraseña
              </label>
              <input
                type="password"
                value={password}
                onChange={e => setPassword(e.target.value)}
                placeholder="••••••••"
                required
                className="w-full bg-[#0D1620] border border-white/10 rounded-lg px-4 py-2.5 text-sm text-[#E8EEF4] placeholder-[#4A5C6A] outline-none focus:border-brand-cyan/50 focus:shadow-[0_0_0_3px_rgba(0,212,255,0.08)] transition-all"
              />
            </div>

            {error && (
              <div className="bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-2.5 text-red-400 text-sm animate-slide-up">
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-brand-cyan text-[#080D12] font-display font-semibold text-sm rounded-lg py-2.5 mt-2 hover:bg-white transition-colors disabled:opacity-50 disabled:cursor-not-allowed shadow-accent"
            >
              {loading ? 'Autenticando...' : 'Iniciar sesión'}
            </button>
          </form>

          {/* Dev hint */}
          <div className="mt-6 pt-5 border-t border-white/5">
            <p className="text-[#4A5C6A] text-xs font-mono mb-2">DEV — credenciales de prueba:</p>
            <div className="space-y-1">
              {[
                ['admin@planetaazul.com', 'admin123', 'Admin'],
                ['chefe@planetaazul.com', 'chefe123', 'Jefe de Área'],
                ['sup@planetaazul.com',   'sup123',   'Supervisor'],
                ['membro@planetaazul.com','membro123','Miembro'],
              ].map(([e, p, label]) => (
                <button
                  key={e}
                  type="button"
                  onClick={() => { setEmail(e); setPassword(p) }}
                  className="w-full text-left px-3 py-1.5 rounded bg-white/3 hover:bg-white/5 text-xs text-[#7D95A8] font-mono transition-colors"
                >
                  <span className="text-brand-cyan">{label}</span> — {e}
                </button>
              ))}
            </div>
          </div>
        </div>

        <p className="text-center text-[#4A5C6A] text-xs mt-6 font-mono">
          v0.1.0 · Planeta Azul Industrial
        </p>
      </div>
    </div>
  )
}
