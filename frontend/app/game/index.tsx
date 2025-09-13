import React from "react";
import GameBoard from "../../components/gameBoard/GameBoard";
import { useUsersQuery } from "@/queries/usePlayersQuery";

const GamePage = () => {
  const { data: users, isLoading, error } = useUsersQuery();
  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error loading players</div>;
  if (!users || users.length === 0) return <div>No players found</div>;
  return (
    <div className="flex justify-center items-center h-screen">
      <GameBoard names={users} />
    </div>
  );
};

export default GamePage;
