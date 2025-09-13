import { useQuery } from "@tanstack/react-query";
import { fetchAllUsers } from "../services/users";

export const useUsersQuery = () => {
  return useQuery({
    queryKey: ["allUsers"], // unique key for this query
    queryFn: fetchAllUsers, // function that returns a Promise
  });
};
