'use client'

import { GuessComparison } from '@/lib/game'
import HintTile from './HintTile'

interface GuessRowProps {
  comparison: GuessComparison
  rowIndex: number
}

const TILE_STAGGER = 150
const ROW_STAGGER = 100

export default function GuessRow({ comparison, rowIndex }: GuessRowProps) {
  const { driver, age, team, nationality, age_hint, debut, wins } = comparison
  const baseDelay = rowIndex * ROW_STAGGER

  const tiles = [
    { label: 'Driver', value: driver.name.split(' ').slice(-1)[0], result: 'red' as const, direction: undefined },
    { label: 'Flag', value: driver.nationality, result: nationality.result, direction: undefined },
    { label: 'Team', value: driver.team, result: team.result, direction: undefined },
    { label: 'Age', value: age, result: age_hint.result, direction: age_hint.direction },
    { label: 'Debut', value: driver.debut_year, result: debut.result, direction: debut.direction },
    { label: 'Wins', value: driver.wins, result: wins.result, direction: wins.direction },
  ]

  return (
    <div className="grid gap-2 w-full" style={{ gridTemplateColumns: 'repeat(6, 1fr)' }}>
      {tiles.map((tile, i) => (
        <HintTile
          key={tile.label}
          label={tile.label}
          value={tile.value}
          result={tile.result}
          direction={tile.direction}
          delay={baseDelay + i * TILE_STAGGER}
        />
      ))}
    </div>
  )
}
