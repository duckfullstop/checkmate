# 🃏 Checkmate.

> If we hit that bullseye, the rest of the dominoes will fall like a house of cards. Checkmate.
>
> -_Zapp Brannigan_

### 📝 A Take-Home Test Exercise for the BBC

This is my solution to a technical assessment I've completed for the BBC. The assessment calls for software that can
correctly score a full game of Blackjack (also known as Pontoon), but I've gone a little further than that.

Yes, the name is exorbitantly stupid. I couldn't think up anything smarter when starting the project.

### 🤩 Features

This software:

* properly handles logic for and scores a game of Blackjack, per the original brief
* has a simple CLI program to simulate a game, called _localjack_
    * which deals from a single deck of 52 cards (configurable)
* exposes basic libraries for building card games, including the concept of a "deck of cards"
  * these libraries are safe to use in threaded, asynchronous environments
* has unit testing for all of the above.
  * Coverage: 100% for all packages, _localjack_ not tested (it's scrappy)
  * Includes defined test cases for the scenarios in the initial brief (`TestBrief*`)

Everything is written in native _Go_ with no external libraries.

> So I really think the thing to take away from all this is that everything is basically fine. So that's all good.
>
> -_Ian Fletcher, W1A_

For more information on why things are laid out as they are, see individual `README.md`'s in each package.
Further discussion points are also labelled as comments that start with `SWEng:`.

## 🤔 Usage Instructions

A makefile is included to make this as simple as possible. To run:

0. Ensure you have `make` and `go` installed (this one's up to you)
1. Move everything here to `$GOROOT/src/github.com/duckfullstop/checkmate`
2. (optionally) Run tests: `make test`
3. Build the app: `make build`

* This drops a binary at `./localjack`

4. Run the game: `./localjack`

## 📄 License

The contents of this repository are licensed under the [MIT license](LICENSE.md). Basically, go wild - if it helps you learn, all the better!