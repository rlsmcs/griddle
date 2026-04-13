'use client'

const HEADERS = ['Driver', 'Flag', 'Team', 'Age', 'Debut', 'Wins']

export default function ColumnHeaders() {
  return (
    <div className="grid gap-2 w-full mb-1" style={{ gridTemplateColumns: 'repeat(6, 1fr)' }}>
      {HEADERS.map(h => (
        <div
          key={h}
          className="text-center py-1.5"
          style={{
            fontSize: '11px',
            fontWeight: 600,
            color: '#aeaeb2',
            letterSpacing: '0.03em',
            fontFamily: '"DM Sans", sans-serif',
          }}
        >
          {h}
        </div>
      ))}
    </div>
  )
}
