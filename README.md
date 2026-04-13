project:
name: griddle
description: formula 1 guess the driver minigame
status: deployed
url: https://griddle-app.vercel.app/

tech_stack:
frontend:
- next.js
- tailwind css
- typescript
backend:
- golang
- rest api

structure:
root:
- app:
- globals.css
- layout.tsx
- page.tsx
- components:
- ColumnHeaders.tsx
- GameOverModal.tsx
- GuessRow.tsx
- HintTile.tsx
- Legend.tsx
- SearchInput.tsx
- data:
- drivers.json
- drivers_sample.json
- public:
- data:
- drivers.json
- lib:
- game.ts
- types.ts
- scripts: []
- config:
- next.config.js
- postcss.config.js
- tailwind.config.js
- tsconfig.json
- next-env.d.ts
- other:
- package.json
- package-lock.json
- README.md
- LICENSE
- .gitignore
- architecture.yaml

deployment:
platform: vercel
type: frontend hosting
backend: external go service just to update the driver dataset with latest data

notes:

* static driver dataset generated via golang script
* dataset updated on a race-ly basis with latest race statistics through the same golang script (i can also add drivers, currently its just 2018 onwards)
* game uses attribute comparison (nationality, team, age, debut, wins)
* six guess limit per session
* simple color and directional hint system (a more logical colour-set will be implemented later , maybe purple for greater than expected ans, yellow for lesser and so on )
* ui is minimal and will be improved later (it is not to scale right now as well)

