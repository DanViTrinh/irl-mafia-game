import AsyncStorage from "@react-native-async-storage/async-storage";
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

export const loginUser = async (
  username: string,
  password: string
): Promise<{ success: boolean; user?: any; error?: string }> => {
  try {
    const { data } = await api.post("/login", { username, password });

    if (data?.token) {
      await AsyncStorage.setItem("jwt", data.token);
      return { success: true, user: { id: data.id, username: data.username } };
    }

    return { success: false, error: "No token returned" };
  } catch (err: any) {
    return {
      success: false,
      error: err.response?.data?.error || "Login failed",
    };
  }
};
