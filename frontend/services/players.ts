import axios from "axios";

export const fetchPlayers = async (): Promise<string[]> => {
  const { data } = await axios.get(`https://randomuser.me/api/?results=${5}`);
  return data.results.map((user: any) => user.name.first);
};
