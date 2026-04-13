'use client'

import { useState, useRef, useEffect } from 'react'
import { Driver } from '@/lib/types'

interface SearchInputProps {
  drivers: Driver[]
  guessedIds: number[]
  onGuess: (driver: Driver) => void
  disabled: boolean
}

export default function SearchInput({ drivers, guessedIds, onGuess, disabled }: SearchInputProps) {
  const [query, setQuery] = useState('')
  const [open, setOpen] = useState(false)
  const [highlighted, setHighlighted] = useState(0)
  const inputRef = useRef<HTMLInputElement>(null)

  const filtered = query.trim().length > 0
    ? drivers
      .filter(d => !guessedIds.includes(d.id) && d.name.toLowerCase().includes(query.toLowerCase()))
      .slice(0, 8)
    : []

  useEffect(() => { setHighlighted(0) }, [query])

  const selectDriver = (driver: Driver) => {
    onGuess(driver)
    setQuery('')
    setOpen(false)
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!open || filtered.length === 0) return
    if (e.key === 'ArrowDown') { e.preventDefault(); setHighlighted(h => Math.min(h + 1, filtered.length - 1)) }
    else if (e.key === 'ArrowUp') { e.preventDefault(); setHighlighted(h => Math.max(h - 1, 0)) }
    else if (e.key === 'Enter') { e.preventDefault(); if (filtered[highlighted]) selectDriver(filtered[highlighted]) }
    else if (e.key === 'Escape') setOpen(false)
  }

  return (
    <div className="relative w-full max-w-sm mx-auto">
      <input
        ref={inputRef}
        type="text"
        value={query}
        onChange={e => { setQuery(e.target.value); setOpen(true) }}
        onFocus={() => setOpen(true)}
        onBlur={() => setTimeout(() => setOpen(false), 150)}
        onKeyDown={handleKeyDown}
        disabled={disabled}
        placeholder={disabled ? 'Game over' : 'Search a driver...'}
        style={{
          width: '100%',
          padding: '12px 16px',
          background: '#ffffff',
          border: '1.5px solid #e0e0e5',
          borderRadius: '12px',
          fontSize: '14px',
          fontFamily: '"DM Sans", sans-serif',
          color: '#1c1c1e',
          outline: 'none',
          transition: 'border-color 0.15s, box-shadow 0.15s',
          boxShadow: '0 1px 3px rgba(0,0,0,0.06)',
        }}
        onFocusCapture={e => {
          e.currentTarget.style.borderColor = '#34c759'
          e.currentTarget.style.boxShadow = '0 0 0 3px rgba(52,199,89,0.12)'
        }}
        onBlurCapture={e => {
          e.currentTarget.style.borderColor = '#e0e0e5'
          e.currentTarget.style.boxShadow = '0 1px 3px rgba(0,0,0,0.06)'
        }}
        autoComplete="off"
        spellCheck={false}
      />

      {open && filtered.length > 0 && (
        <div
          className="absolute left-0 right-0 z-50 mt-1 fade-in"
          style={{
            background: '#ffffff',
            border: '1.5px solid #e0e0e5',
            borderRadius: '12px',
            boxShadow: '0 8px 24px rgba(0,0,0,0.1)',
            overflow: 'hidden',
          }}
        >
          {filtered.map((driver, i) => (
            <button
              key={driver.id}
              onMouseDown={() => selectDriver(driver)}
              onMouseEnter={() => setHighlighted(i)}
              style={{
                width: '100%',
                textAlign: 'left',
                padding: '10px 16px',
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                background: i === highlighted ? '#f2f2f7' : 'transparent',
                border: 'none',
                borderBottom: i < filtered.length - 1 ? '1px solid #f2f2f7' : 'none',
                fontFamily: '"DM Sans", sans-serif',
                fontSize: '13px',
                color: '#1c1c1e',
                cursor: 'pointer',
                transition: 'background 0.1s',
              }}
            >
              <span style={{ fontWeight: 500 }}>{driver.name}</span>
              <span style={{ color: '#aeaeb2', fontSize: '11px' }}>{driver.team}</span>
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
