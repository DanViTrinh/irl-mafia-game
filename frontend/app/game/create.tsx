"use client";

import { useState } from "react";

export default function CreateGamePage() {
  const [playerIds, setPlayerIds] = useState<string[]>([""]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handlePlayerChange = (index: number, value: string) => {
    const newPlayers = [...playerIds];
    newPlayers[index] = value;
    setPlayerIds(newPlayers);
  };

  const addPlayer = () => setPlayerIds([...playerIds, ""]);
  const removePlayer = (index: number) =>
    setPlayerIds(playerIds.filter((_, i) => i !== index));

  const createGame = async () => {
    setLoading(true);
    setError("");
    try {
      const res = await fetch("/api/games", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ playerIds }),
      });
      if (!res.ok) throw new Error("Failed to create game");
      const data = await res.json();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-xl mx-auto mt-12 p-4 bg-white rounded shadow">
      <h1 className="text-2xl font-bold mb-4">Create New Game</h1>

      {playerIds.map((player, index) => (
        <div key={index} className="flex items-center mb-2">
          <input
            type="text"
            placeholder={`Player ${index + 1} ID`}
            value={player}
            onChange={(e) => handlePlayerChange(index, e.target.value)}
            className="flex-1 border rounded p-2 mr-2"
          />
          {playerIds.length > 1 && (
            <button
              onClick={() => removePlayer(index)}
              className="text-red-500 font-bold"
            >
              X
            </button>
          )}
        </div>
      ))}

      <button
        onClick={addPlayer}
        className="mb-4 px-4 py-2 bg-blue-500 text-white rounded"
      >
        Add Player
      </button>

      {error && <p className="text-red-500 mb-2">{error}</p>}

      <button
        onClick={createGame}
        disabled={loading}
        className={`w-full px-4 py-2 text-white rounded ${loading ? "bg-gray-400" : "bg-green-500"}`}
      >
        {loading ? "Creating..." : "Create Game"}
      </button>
    </div>
  );
}
