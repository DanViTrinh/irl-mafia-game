import api from "./api";

export const fetchAllUsers = async (): Promise<string[]> => {
  const { data } = await api.get("/users");
  return data.results.map((user: any) => user.name);
};

export const signupUser = async (
  username: string,
  password: string
): Promise<any> => {
  const { data } = await api.post("/signup", { username, password });
  return data;
};
