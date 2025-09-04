import { useQuery } from "@tanstack/react-query";
import { fetchPlayers } from "../services/players";

export const usePlayersQuery = () => {
  return useQuery({
    queryKey: ["players"], // unique key for this query
    queryFn: fetchPlayers, // function that returns a Promise
  });
};
