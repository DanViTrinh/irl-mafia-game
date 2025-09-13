import React from "react";
import GameBoard from "../../components/gameBoard/GameBoard";
import { useUsersQuery } from "@/queries/useUsersQuery";
import { useLocalSearchParams } from "expo-router/build/hooks";

const GameScreen = () => {
  const { id } = useLocalSearchParams();
  const { data: users, isLoading, error } = useUsersQuery();
  if (isLoading) return <div>Loading...</div>;
  if (error) {
    console.log("Error fetching users:", error);
    return <div>Error loading users</div>;
  }
  if (!users || users.length === 0) return <div>No users found</div>;
  return (
    <div className="flex justify-center items-center h-screen">
      <h1>Game ID: {id}</h1>
      <GameBoard names={users} />
    </div>
  );
};

export default GameScreen;
