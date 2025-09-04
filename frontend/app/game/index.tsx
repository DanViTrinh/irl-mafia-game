import React from "react";
import GameBoard from "../../components/gameBoard/GameBoard";

const GamePage = () => {
  return (
    <div className="flex justify-center items-center h-screen">
      <GameBoard names={["ole", "petter", "dan", "asdf", "asldfkj"]} />
    </div>
  );
};

export default GamePage;
