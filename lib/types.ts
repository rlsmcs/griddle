export interface Driver {
  id: number
  code: string
  name: string
  nationality: string
  date_of_birth: string
  team: string
  past_teams: string[]
  car_number: number
  debut_year: number
  wins: number
}

export type HintResult = 'green' | 'yellow' | 'half-green' | 'red'

export interface GuessResult {
  driver: Driver
  hints: {
    name: HintResult
    team: HintResult
    nationality: HintResult
    debut_year: HintResult & { direction?: 'up' | 'down' | 'exact' }
    wins: HintResult & { direction?: 'up' | 'down' | 'exact' }
  }
  teamHint: HintResult
  nationalityHint: HintResult
  debutHint: { result: HintResult; direction: 'up' | 'down' | 'exact' }
  winsHint: { result: HintResult; direction: 'up' | 'down' | 'exact' }
}
