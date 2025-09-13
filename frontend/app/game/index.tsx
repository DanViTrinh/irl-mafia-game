import { useRouter } from "expo-router";
import React, { useState } from "react";

export default function GameIndex() {
  const router = useRouter();
  const [showJoin, setShowJoin] = useState(false);
  const [showCreate, setShowCreate] = useState(false);
  const [gameCode, setGameCode] = useState("");
  const [newGameName, setNewGameName] = useState("");

  const joinGame = () => {
    if (gameCode) {
      router.push(`/game/${gameCode}`);
    } else {
      alert("Please enter game code to join a game.");
    }
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-100">
      <div className="bg-white shadow-lg rounded-lg p-8 w-full max-w-md">
        <h1 className="text-3xl font-bold mb-6 text-center">IRL Mafia Game</h1>
        <div className="flex flex-col gap-4">
          <button
            className="bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition"
            onClick={() => {
              setShowJoin(true);
              setShowCreate(false);
            }}
          >
            Join Game
          </button>
          <button
            className="bg-green-600 text-white py-2 px-4 rounded hover:bg-green-700 transition"
            onClick={() => {
              setShowCreate(true);
              setShowJoin(false);
            }}
          >
            Create Game
          </button>
        </div>

        {showJoin && (
          <form className="mt-6 flex flex-col gap-3">
            <input
              type="text"
              placeholder="Game Code"
              value={gameCode}
              onChange={(e) => setGameCode(e.target.value)}
              className="border rounded px-3 py-2"
            />
            <button
              type="submit"
              className="bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700 transition"
              onClick={(e) => {
                joinGame();
              }}
            >
              Join
            </button>
          </form>
        )}

        {showCreate && (
          <form className="mt-6 flex flex-col gap-3">
            <input
              type="text"
              placeholder="Game Name"
              value={newGameName}
              onChange={(e) => setNewGameName(e.target.value)}
              className="border rounded px-3 py-2"
            />
            <button
              type="submit"
              className="bg-green-600 text-white py-2 px-4 rounded hover:bg-green-700 transition"
              onClick={(e) => {
                e.preventDefault();
                // TODO: handle create game logic
              }}
            >
              Create
            </button>
          </form>
        )}
      </div>
    </div>
  );
}
