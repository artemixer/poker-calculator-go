# poker-calculator-go
**poker-calculator-go** is a tool for estimating the odds of winning/getting certain combinations given a state of game in Texas Hold'em Poker. 


## Installation
Grab the binary for your respective system from the releases section of this repo
<br/>Or, you could build it yourself: make sure you have Golang installed by running ```go version```
```
git clone https://github.com/artemixer/poker-calculator-go
cd poker-calculator-go
go build poker-calculator.go
```
  
## Usage
```
./poker-calculator -i input_state.json -iter 50000 -verbose false
```
Output:
<br/><br/><img width="305" alt="Screenshot" src="https://github.com/user-attachments/assets/de4b2e10-fecf-42ee-8ca8-bbfed397ae5b" />


The tool reads the current game state from the input file (```-i``` parameter), a sample file is included in the repository. The cards in the file are
contained in the format where ```Ah``` stands for Ace of Hearts and ```Tc``` is the 10 of Clubs. The community cards not yet revealed should be replaced with
```not_drawn``` and will be simulated randomly.


## Methodology
Given that a game of poker with just 4 players can have up to billions of possible combinations, this tool uses the Monte Carlo method to solve this issue. 
Meaning instead of simulating every possible outcome in a given game, we generate random games given the set parameters, 
which tends to estimate the real odds with sufficient accuracy. For example at 20000 iterations it is accurate in the range of 1%, at 50000 iterations
0.5%. 

<br/>
<br/>
<b>Feel free to open pull requests and issues, I will do my best to resolve or at least answer all of them</b>
