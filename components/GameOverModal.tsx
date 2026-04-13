'use client'

import { Driver } from '@/lib/types'

interface GameOverModalProps {
  won: boolean
  target: Driver
  attempts: number
  onClose: () => void
}

export default function GameOverModal({ won, target, attempts, onClose }: GameOverModalProps) {
  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4"
      style={{ background: 'rgba(0,0,0,0.3)', backdropFilter: 'blur(8px)' }}
      onClick={onClose}
    >
      <div
        className="fade-in relative max-w-sm w-full p-8"
        style={{
          background: '#ffffff',
          borderRadius: '20px',
          boxShadow: '0 24px 60px rgba(0,0,0,0.15)',
          border: '1.5px solid #e0e0e5',
        }}
        onClick={e => e.stopPropagation()}
      >
        {/* Status pill */}
        <div className="flex justify-center mb-6">
          <div
            style={{
              padding: '6px 16px',
              borderRadius: '999px',
              background: won ? '#e8f8ed' : '#fff0ef',
              border: `1.5px solid ${won ? '#34c759' : '#ff3b30'}`,
              fontSize: '12px',
              fontWeight: 600,
              color: won ? '#1c6b32' : '#7a1a15',
              fontFamily: '"DM Sans", sans-serif',
              letterSpacing: '0.02em',
            }}
          >
            {won ? `Solved in ${attempts} ${attempts === 1 ? 'guess' : 'guesses'}` : 'Out of guesses'}
          </div>
        </div>

        <div style={{ textAlign: 'center', marginBottom: '24px' }}>
          <div style={{ fontSize: '12px', color: '#aeaeb2', marginBottom: '8px', fontFamily: '"DM Sans"' }}>
            The driver was
          </div>
          <div style={{ fontSize: '22px', fontWeight: 700, color: '#1c1c1e', fontFamily: '"DM Sans"', lineHeight: 1.2 }}>
            {target.name}
          </div>
        </div>

        <div
          className="grid grid-cols-2 gap-3 mb-6"
          style={{ fontFamily: '"DM Sans", sans-serif', fontSize: '13px' }}
        >
          {[
            { label: 'Team', value: target.team },
            { label: 'Nationality', value: target.nationality },
            { label: 'Debut', value: target.debut_year },
            { label: 'Wins', value: target.wins },
          ].map(row => (
            <div
              key={row.label}
              style={{
                background: '#f7f7f7',
                borderRadius: '10px',
                padding: '10px 12px',
              }}
            >
              <div style={{ fontSize: '10px', color: '#aeaeb2', fontWeight: 600, marginBottom: '2px', textTransform: 'uppercase', letterSpacing: '0.04em' }}>
                {row.label}
              </div>
              <div style={{ fontWeight: 600, color: '#1c1c1e' }}>{row.value}</div>
            </div>
          ))}
        </div>

        <button
          onClick={onClose}
          style={{
            width: '100%',
            padding: '12px',
            borderRadius: '12px',
            background: '#1c1c1e',
            border: 'none',
            color: '#ffffff',
            fontSize: '14px',
            fontWeight: 600,
            fontFamily: '"DM Sans", sans-serif',
            cursor: 'pointer',
          }}
        >
          Close
        </button>
      </div>
    </div>
  )
}
