import { Text, View } from "react-native";
import "../global.css";
import { Button } from "@react-navigation/elements";
import { router } from "expo-router";
import LoginScreen from "./login";

export default function Index() {
  return LoginScreen();
}
