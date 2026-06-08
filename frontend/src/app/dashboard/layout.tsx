'use client'

import { useEffect } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import Link from 'next/link'
import { useAuth } from '@/hooks/useAuth'
import {
  LayoutDashboard, ListTodo, ShieldCheck, BarChart2,
  FileText, ScanSearch, LogOut, Bell, Settings
} from 'lucide-react'
import { ROLE_LABELS } from '@/types'
import { Logo } from '@/components/ui/Logo'

const NAV = [
  { href: '/dashboard',  icon: LayoutDashboard, label: 'Início',       key: 'dashboard' },
  { href: '/tasks',      icon: ListTodo,        label: 'Tasks',         key: 'tasks' },
  { href: '/access',     icon: ShieldCheck,     label: 'Acesso',        key: 'access' },
  { href: '/data',       icon: BarChart2,       label: 'Dados / KPIs',  key: 'data' },
  { href: '/forms',      icon: FileText,        label: 'Formulários',   key: 'forms' },
  { href: '/ocr',        icon: ScanSearch,      label: 'Documentos',    key: 'ocr' },
]

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const { user, loading, logout } = useAuth()
  const router = useRouter()
  const pathname = usePathname()

  useEffect(() => {
    if (!loading && !user) router.push('/login')
  }, [user, loading, router])

  if (loading || !user) {
    return (
      <div className="min-h-screen bg-[#080D12] flex items-center justify-center">
        <div className="w-6 h-6 border-2 border-brand-cyan/40 border-t-brand-cyan rounded-full animate-spin" />
      </div>
    )
  }

  const initials = user.name.split(' ').map(n => n[0]).join('').slice(0, 2).toUpperCase()

  return (
    <div className="min-h-screen bg-[#080D12] flex">
      {/* Sidebar */}
      <aside className="w-[220px] shrink-0 border-r border-white/[0.06] flex flex-col bg-[#080D12]">
        {/* Logo */}
        <div className="h-16 flex items-center px-5 border-b border-white/[0.06]">
          <Logo size={32} className="mr-3 rounded-full" />
          <span className="font-display font-bold text-sm text-[#E8EEF4]">Planeta Azul</span>
        </div>

        {/* Nav */}
        <nav className="flex-1 px-3 py-4 space-y-0.5 overflow-y-auto">
          {NAV.map(({ href, icon: Icon, label }) => {
            const active = pathname === href || (href !== '/dashboard' && pathname.startsWith(href))
            return (
              <Link
                key={href}
                href={href}
                className={`flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all group ${
                  active
                    ? 'bg-brand-cyan/10 text-brand-cyan'
                    : 'text-[#7D95A8] hover:text-[#E8EEF4] hover:bg-white/[0.04]'
                }`}
              >
                <Icon size={16} className={active ? 'text-brand-cyan' : 'text-[#4A5C6A] group-hover:text-[#7D95A8]'} />
                <span className={active ? 'font-medium' : ''}>{label}</span>
                {active && <div className="ml-auto w-1 h-1 rounded-full bg-brand-cyan" />}
              </Link>
            )
          })}
        </nav>

        {/* Bottom: user + settings */}
        <div className="px-3 pb-4 border-t border-white/[0.06] pt-3 space-y-0.5">
          <Link href="/settings" className="flex items-center gap-3 px-3 py-2 rounded-lg text-[#7D95A8] hover:text-[#E8EEF4] hover:bg-white/[0.04] text-sm transition-all">
            <Settings size={15} />
            Configurações
          </Link>
          <button
            onClick={logout}
            className="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-[#7D95A8] hover:text-red-400 hover:bg-red-500/[0.06] text-sm transition-all"
          >
            <LogOut size={15} />
            Sair
          </button>
          {/* User pill */}
          <div className="flex items-center gap-2.5 px-3 py-2.5 mt-1 rounded-lg bg-white/[0.03] border border-white/[0.05]">
            <div className="w-7 h-7 rounded-full bg-brand-cyan/20 border border-brand-cyan/30 flex items-center justify-center text-brand-cyan text-xs font-display font-bold shrink-0">
              {initials}
            </div>
            <div className="min-w-0">
              <p className="text-xs text-[#E8EEF4] font-medium truncate">{user.name.split(' ')[0]}</p>
              <p className="text-[10px] text-[#4A5C6A] font-mono truncate">{ROLE_LABELS[user.role]}</p>
            </div>
          </div>
        </div>
      </aside>

      {/* Main content */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Topbar */}
        <header className="h-16 border-b border-white/[0.06] flex items-center justify-between px-6 shrink-0 bg-[#080D12]">
          <div>
            <h1 className="font-display font-semibold text-base text-[#E8EEF4]">
              {NAV.find(n => pathname === n.href || pathname.startsWith(n.href + '/'))?.label ?? 'Início'}
            </h1>
            {user.area_id && (
              <p className="text-[#4A5C6A] text-xs font-mono mt-0.5">Área: {user.area_id.slice(0, 8)}</p>
            )}
          </div>
          <div className="flex items-center gap-3">
            <button className="relative p-2 rounded-lg text-[#7D95A8] hover:text-[#E8EEF4] hover:bg-white/[0.04] transition-all">
              <Bell size={16} />
              <span className="absolute top-1.5 right-1.5 w-1.5 h-1.5 rounded-full bg-brand-cyan" />
            </button>
          </div>
        </header>

        {/* Page content */}
        <main className="flex-1 overflow-auto p-6">
          {children}
        </main>
      </div>
    </div>
  )
}
