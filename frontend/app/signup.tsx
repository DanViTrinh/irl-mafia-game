import React, { useState } from "react";
import { View, Text, TextInput, TouchableOpacity, Alert } from "react-native";
import { useRouter } from "expo-router";
import {
  QueryClient,
  QueryClientProvider,
  useMutation,
} from "@tanstack/react-query";
import api from "../services/api";
import { signupUser } from "@/services/users";

export default function SignupScreen() {
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const signupMutation = useMutation({
    mutationFn: async () => {
      const { data } = await signupUser(username, password);
      return data;
    },
    onSuccess: (data) => {
      Alert.alert("Success", "Account created successfully");
      console.log("Signup response data:", data);
      router.push("/login");
    },
    onError: (error: any) => {
      Alert.alert("Error", error.response?.data?.message || "Signup failed");
    },
  });

  const handleSignup = () => {
    if (!username || !password) {
      Alert.alert("Error", "Please enter username and password");
      return;
    }

    signupMutation.mutate();
  };

  const handleLoginRedirect = () => router.push("/login");

  return (
    <View className="flex-1 justify-center items-center bg-gray-100 px-6">
      <Text className="text-3xl font-bold mb-8">Sign Up</Text>

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
        onPress={handleSignup}
        className="w-full p-3 mb-4 bg-green-500 rounded"
      >
        <Text className="text-center text-white font-semibold">Sign Up</Text>
      </TouchableOpacity>

      <TouchableOpacity onPress={handleLoginRedirect}>
        <Text className="text-blue-500">Already have an account? Login</Text>
      </TouchableOpacity>
    </View>
  );
}
