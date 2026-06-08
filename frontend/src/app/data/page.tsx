'use client'
export default function DataPage() {
  return <Placeholder label="Datos y KPIs" desc="Módulo 3 — KPIs y exportaciones de la base de datos de la fábrica" />
}
function Placeholder({ label, desc }: { label: string; desc: string }) {
  return (
    <div className="flex items-center justify-center py-24 animate-fade-in">
      <div className="text-center">
        <div className="w-12 h-12 rounded-full bg-brand-cyan/10 border border-brand-cyan/20 flex items-center justify-center mx-auto mb-4">
          <div className="w-2 h-2 rounded-full bg-brand-cyan" /></div>
        <h2 className="font-display text-xl font-bold text-[#E8EEF4] mb-2">{label}</h2>
        <p className="text-[#4A5C6A] text-sm font-mono">{desc}</p>
      </div>
    </div>
  )
}
