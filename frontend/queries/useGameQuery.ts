import { useQuery } from "@tanstack/react-query";
import { fetchUsernameFromGameId } from "@/services/game";

export const useGameQuery = (gameId: string) => {
  return useQuery({
    queryKey: ["allUsers", gameId], // include the param
    queryFn: () => fetchUsernameFromGameId(gameId), // pass a function
    enabled: !!gameId, // optional, prevents running with empty id
  });
};
