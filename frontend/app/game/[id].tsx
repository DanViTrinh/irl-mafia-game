import React from "react";
import GameBoard from "../../components/gameBoard/GameBoard";
import { useLocalSearchParams } from "expo-router/build/hooks";
import { useGameQuery } from "@/queries/useGameQuery";

const GameScreen = () => {
  const { id } = useLocalSearchParams();
  const { data: users, isLoading, error } = useGameQuery(id as string);
  if (isLoading) return <div>Loading...</div>;
  if (error) {
    console.log("Error fetching users:", error);
    return <div>Error loading users</div>;
  }
  if (!users || users.length === 0) return <div>No users found</div>;
  return (
    <div className="flex justify-center items-center h-screen">
      <GameBoard names={users} />
    </div>
  );
};

export default GameScreen;
