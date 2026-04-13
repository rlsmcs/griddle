'use client'

const items = [
  { color: '#34c759', bg: '#e8f8ed', label: 'Exact match' },
  { color: '#ff9f0a', bg: '#fff5e6', label: 'Close / past team' },
  { color: '#ff3b30', bg: '#fff0ef', label: 'No match' },
]

export default function Legend() {
  return (
    <div className="flex flex-wrap gap-x-5 gap-y-2 justify-center">
      {items.map(item => (
        <div key={item.label} className="flex items-center gap-2">
          <div
            style={{
              width: '14px',
              height: '14px',
              borderRadius: '4px',
              background: item.bg,
              border: `1.5px solid ${item.color}`,
              flexShrink: 0,
            }}
          />
          <span style={{ fontSize: '11px', color: '#6c6c70', fontFamily: '"DM Sans", sans-serif' }}>
            {item.label}
          </span>
        </div>
      ))}
    </div>
  )
}
