import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'Griddle',
  description: 'Guess the Formula1 driver in 6 tries. New driver every day.',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        {children}
      </body>
    </html>
  )
}
