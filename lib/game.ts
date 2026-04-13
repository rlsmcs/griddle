import { Driver } from './types'

export type HintResult = 'green' | 'half-green' | 'yellow' | 'red'

export interface DirectionalHint {
  result: HintResult
  direction: 'up' | 'down' | 'exact'
}

export interface GuessComparison {
  driver: Driver
  age: number
  team: { result: HintResult }
  nationality: { result: HintResult }
  age_hint: DirectionalHint
  debut: DirectionalHint
  wins: DirectionalHint
}

export function getDriverAge(driver: Driver): number {
  const dob = new Date(driver.date_of_birth)
  const today = new Date()
  let age = today.getFullYear() - dob.getFullYear()
  const m = today.getMonth() - dob.getMonth()
  if (m < 0 || (m === 0 && today.getDate() < dob.getDate())) age--
  return age
}

export function getDailyDriver(drivers: Driver[]): Driver {
  const epoch = new Date('2024-01-01').getTime()
  const now = new Date().getTime()
  const daysSinceEpoch = Math.floor((now - epoch) / (1000 * 60 * 60 * 24))
  const index = daysSinceEpoch % drivers.length
  return drivers[index]
}

export function compareGuess(guess: Driver, target: Driver): GuessComparison {
  //green:  guess team == target team (exact)
  //yellow: guess team appears in target's past teams
  //red:  no connection
  const guessTeam = guess.team
  const targetTeam = target.team
  const targetPastTeams = target.past_teams ?? []

  let teamResult: HintResult
  if (guessTeam === targetTeam) {
    teamResult = 'green'
  } else if (targetPastTeams.includes(guessTeam)) {
    teamResult = 'yellow'
  } else {
    teamResult = 'red'
  }

  // nationality( make it flag later phase but for now ive jsut kept it in words)
  const nationalityResult: HintResult =
    guess.nationality === target.nationality ? 'green' : 'red'

  //age
  const guessAge = getDriverAge(guess)
  const targetAge = getDriverAge(target)
  const ageDiff = guessAge - targetAge
  let ageResult: HintResult
  if (ageDiff === 0) ageResult = 'green'
  else if (Math.abs(ageDiff) <= 2) ageResult = 'yellow'
  else ageResult = 'red'
  const ageDirection = ageDiff === 0 ? 'exact' : ageDiff > 0 ? 'down' : 'up'

  //debut
  const debutDiff = guess.debut_year - target.debut_year
  let debutResult: HintResult
  if (debutDiff === 0) debutResult = 'green'
  else if (Math.abs(debutDiff) <= 2) debutResult = 'yellow'
  else debutResult = 'red'
  const debutDirection = debutDiff === 0 ? 'exact' : debutDiff > 0 ? 'down' : 'up'

  //GP wins
  const winsDiff = guess.wins - target.wins
  let winsResult: HintResult
  if (winsDiff === 0) winsResult = 'green'
  else if (Math.abs(winsDiff) <= 5) winsResult = 'yellow'
  else winsResult = 'red'
  const winsDirection = winsDiff === 0 ? 'exact' : winsDiff > 0 ? 'down' : 'up'

  return {
    driver: guess,
    age: guessAge,
    team: { result: teamResult },
    nationality: { result: nationalityResult },
    age_hint: { result: ageResult, direction: ageDirection },
    debut: { result: debutResult, direction: debutDirection },
    wins: { result: winsResult, direction: winsDirection },
  }
}
