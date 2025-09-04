import React from "react";
import GameBoard from "../../components/gameBoard/GameBoard";
import { usePlayersQuery } from "@/queries/usePlayersQuery";

const GamePage = () => {
  const { data: players, isLoading, error } = usePlayersQuery();
  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error loading players</div>;
  if (!players || players.length === 0) return <div>No players found</div>;
  return (
    <div className="flex justify-center items-center h-screen">
      <GameBoard names={players} />
    </div>
  );
};

export default GamePage;
