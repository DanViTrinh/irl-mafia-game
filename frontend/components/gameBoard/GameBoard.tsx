import React, { useState } from "react";
import { View } from "react-native";
import Cell from "./Cell";

type GameBoardProps = {
  names: string[];
};

// Fisher–Yates shuffle
const shuffleArray = (arr: string[]) => {
  const array = [...arr];
  for (let i = array.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [array[i], array[j]] = [array[j], array[i]];
  }
  return array;
};

const GameBoard = ({ names }: GameBoardProps) => {
  // Compute the minimum size needed to fit all players
  const size = Math.ceil(Math.sqrt(names.length));
  const totalCells = size * size;

  const [board] = useState(() => {
    let shuffledNames = shuffleArray(names);

    // Repeat names to fill board if needed
    while (shuffledNames.length < totalCells) {
      shuffledNames = shuffledNames.concat(shuffleArray(names));
    }

    return shuffledNames.slice(0, totalCells);
  });

  return (
    <View style={{ flexDirection: "row", flexWrap: "wrap", width: size * 64 }}>
      {board.map((name, i) => (
        <Cell key={i} name={name} marked={false} onPress={() => {}} />
      ))}
    </View>
  );
};

export default GameBoard;
