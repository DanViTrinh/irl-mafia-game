import axios from "axios";
import AsyncStorage from "@react-native-async-storage/async-storage";

// Create instance
const api = axios.create({
  baseURL: "http://localhost:8080/", // Replace with your API base URL
  headers: {
    "Content-Type": "application/json",
  },
});

// Optional: add JWT automatically to every request
api.interceptors.request.use(
  (config) => {
    const token = AsyncStorage.getItem("jwt");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

export default api;
