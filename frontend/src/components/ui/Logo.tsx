import Image from 'next/image'

interface LogoProps {
  size?: number
  className?: string
}

export function Logo({ size = 64, className = '' }: LogoProps) {
  return (
    <Image
      src="/brand/logo.png"
      alt="Planeta Azul"
      width={size}
      height={size}
      className={className}
      priority
    />
  )
}
