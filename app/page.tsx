'use client'

import { useState, useEffect } from 'react'
import { Driver } from '@/lib/types'
import { getDailyDriver, compareGuess, GuessComparison } from '@/lib/game'
import SearchInput from '@/components/SearchInput'
import GuessRow from '@/components/GuessRow'
import ColumnHeaders from '@/components/ColumnHeaders'
import GameOverModal from '@/components/GameOverModal'
import Legend from '@/components/Legend'

const MAX_GUESSES = 6

export default function Home() {
  const [drivers, setDrivers] = useState<Driver[]>([])
  const [target, setTarget] = useState<Driver | null>(null)
  const [guesses, setGuesses] = useState<GuessComparison[]>([])
  const [gameState, setGameState] = useState<'playing' | 'won' | 'lost'>('playing')
  const [showModal, setShowModal] = useState(false)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/data/drivers.json')
      .then(r => r.json())
      .then((data: Driver[]) => {
        setDrivers(data)
        setTarget(getDailyDriver(data))
        setLoading(false)
      })
  }, [])

  const handleGuess = (driver: Driver) => {
    if (!target || gameState !== 'playing') return
    if (guesses.some(g => g.driver.id === driver.id)) return

    const comparison = compareGuess(driver, target)
    const newGuesses = [...guesses, comparison]
    setGuesses(newGuesses)

    if (driver.id === target.id) {
      setGameState('won')
      setTimeout(() => setShowModal(true), newGuesses.length * 100 + 6 * 150 + 600)
    } else if (newGuesses.length >= MAX_GUESSES) {
      setGameState('lost')
      setTimeout(() => setShowModal(true), newGuesses.length * 100 + 6 * 150 + 600)
    }
  }

  const guessedIds = guesses.map(g => g.driver.id)
  const remaining = MAX_GUESSES - guesses.length

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div style={{ color: '#aeaeb2', fontSize: '14px', fontFamily: '"DM Sans", sans-serif' }}>
          Loading...
        </div>
      </div>
    )
  }

  return (
    <main
      className="min-h-screen flex flex-col items-center px-4"
      style={{ paddingTop: '40px', paddingBottom: '40px', maxWidth: '700px', margin: '0 auto' }}
    >
      <header className="w-full text-center mb-8">
        <div
          style={{
            fontSize: 'clamp(22px, 5vw, 32px)',
            fontWeight: 700,
            color: '#1c1c1e',
            letterSpacing: '-0.02em',
            fontFamily: '"DM Sans", sans-serif',
          }}
        >
          Griddle
        </div>
        <div
          style={{
            fontSize: '13px',
            color: '#aeaeb2',
            marginTop: '4px',
            fontFamily: '"DM Sans", sans-serif',
          }}
        >
          Guess the F1 driver · {remaining} {remaining === 1 ? 'guess' : 'guesses'} remaining
        </div>

        <div className="flex gap-1.5 justify-center mt-4">
          {Array.from({ length: MAX_GUESSES }).map((_, i) => (
            <div
              key={i}
              style={{
                width: '8px',
                height: '8px',
                borderRadius: '50%',
                background: i < guesses.length
                  ? (gameState === 'won' && i === guesses.length - 1 ? '#34c759' : '#c8c8d0')
                  : '#e0e0e5',
                transition: 'background 0.3s',
              }}
            />
          ))}
        </div>
      </header>

      <div className="w-full">
        <ColumnHeaders />
      </div>

      <div className="w-full mb-3" style={{ height: '1px', background: '#e0e0e5' }} />

      <div className="w-full flex flex-col gap-2 mb-6">
        {guesses.map((comparison, i) => (
          <GuessRow key={comparison.driver.id} comparison={comparison} rowIndex={i} />
        ))}

        {gameState === 'playing' && Array.from({ length: remaining }).map((_, i) => (
          <div
            key={`empty-${i}`}
            className="grid gap-2 w-full"
            style={{ gridTemplateColumns: 'repeat(6, 1fr)' }}
          >
            {Array.from({ length: 6 }).map((_, j) => (
              <div
                key={j}
                style={{
                  aspectRatio: '1 / 1.15',
                  background: '#ffffff',
                  border: '1.5px solid #e0e0e5',
                  borderRadius: '10px',
                  opacity: 1 - i * 0.1,
                }}
              />
            ))}
          </div>
        ))}
      </div>

      <div className="w-full mb-8">
        <SearchInput
          drivers={drivers}
          guessedIds={guessedIds}
          onGuess={handleGuess}
          disabled={gameState !== 'playing'}
        />
      </div>

      <Legend />

      <div
        className="mt-auto pt-8"
        style={{ fontSize: '11px', color: '#c8c8d0', fontFamily: '"DM Sans", sans-serif', textAlign: 'center' }}
      >
        New driver every day · 2018–2026
      </div>

      {showModal && target && (
        <GameOverModal
          won={gameState === 'won'}
          target={target}
          attempts={guesses.length}
          onClose={() => setShowModal(false)}
        />
      )}
    </main>
  )
}
