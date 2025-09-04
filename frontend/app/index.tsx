import { Text, View } from "react-native";
import "../global.css";
import { Button } from "@react-navigation/elements";
import { router } from "expo-router";

export default function Index() {
  return (
    <View className="flex-1 items-center justify-center bg-white space-y-10">
      <Text className="text-xl font-bold text-blue-500">
        Welcome to the FLU game!
      </Text>
      <Button onPress={() => router.push("/game")}>Start</Button>
    </View>
  );
}
