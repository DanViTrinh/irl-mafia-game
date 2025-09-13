import api from "./api";

export const fetchUsernameFromGameId = async (
  gameId: string
): Promise<string[]> => {
  console.log("Fetching users for game ID:", gameId);
  const { data } = await api.get(`/games/${gameId}/players`);
  console.log(data);
  return data;
};
