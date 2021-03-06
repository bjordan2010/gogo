Gogo is my 2nd attempt at a Go (weiqi, baduk) game in Go (lang).  Borrowed some Shiny example code for a good start. I also used *Deep Learning and the Game of Go* by Max Pumperla et al for the game state and AI.

## Todo List
- Fix existing bugs/problems
- [x] Figure out how to use command line options to set board size and scale (19,18 and 9,35) May improve later but it works.
- [x] Change white to right click (3) and remove stone to middle click (2)
- [x] Do more manual testing on 9x9 to understand code
- [ ] Add coordinates
- [ ] Add indicator to last move for each color
- [ ] Add sound for stone placement
- [ ] Upgrade to 1.17
- [ ] Add some written tests to understand shiny code
---
- Add Keyboard commands
- [x] Fix cmd+q not having focus to close even after focusing game
- [x] cmd+n to new game 19x19
- [x] cmd+m to new game 9x9
- [ ] cmd+o to resign as white
- [ ] cmd+p to resign as black
- [ ] cmd+b to pass as black
- [ ] cmd+w to pass as white
---
- Implement Go rules
- [x] Capturing stones
- [ ] Disallow playing out of turn
- [ ] Disallow placing new stone if stone already in place
- [ ] Implement passing
- [ ] Implement resigning
- [ ] Implement scoring after consecutive passes
- [x] Track stones on board using "bunch" (connected stones)
- [x] Track liberties for each stone
- [ ] Track game state part 2
- [ ] Add respecting ko
- [ ] Fix bad captures
---
- Utilities
- [ ] Load board from sgf file
- [ ] Save board to sgf file
- [ ] Create dumb AI
- [ ] Implement review with KataGo
- [ ] Make the AI learn
