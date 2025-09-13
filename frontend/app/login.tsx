import React, { use, useState } from "react";
import { View, Text, TextInput, TouchableOpacity, Alert } from "react-native";
import { useRouter } from "expo-router";
import { useMutation } from "@tanstack/react-query";
import { loginUser } from "@/services/users";

export default function LoginScreen() {
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const loginMutation = useMutation({
    mutationFn: async () => {
      const { success, user } = await loginUser(username, password);
      return { success, user };
    },
    onSuccess: () => {
      Alert.alert("Success", "Logged in successfully");
      router.push("/game");
    },
    onError: (error: any) => {
      Alert.alert("Error", error.response?.data?.message || "Login failed");
    },
  });

  const handleLogin = () => {
    if (!username || !password) {
      Alert.alert("Error", "Please enter username and password");
      return;
    }

    loginMutation.mutate();
  };

  const handleSignup = () => {
    router.push("/signup");
  };

  return (
    <View className="flex-1 justify-center items-center bg-gray-100 px-6">
      <Text className="text-3xl font-bold mb-8">Login</Text>

      <TextInput
        placeholder="Username"
        value={username}
        onChangeText={setUsername}
        className="w-full px-4 py-3 mb-4 bg-white rounded border border-gray-300 h-12"
      />

      <TextInput
        placeholder="Password"
        value={password}
        onChangeText={setPassword}
        secureTextEntry
        className="w-full px-4 py-3 mb-6 bg-white rounded border border-gray-300 h-12"
      />

      <TouchableOpacity
        onPress={handleLogin}
        className="w-full p-3 mb-4 bg-blue-500 rounded"
      >
        <Text className="text-center text-white font-semibold">Login</Text>
      </TouchableOpacity>

      <TouchableOpacity onPress={handleSignup}>
        <Text className="text-blue-500">Don't have an account? Sign up</Text>
      </TouchableOpacity>
    </View>
  );
}
