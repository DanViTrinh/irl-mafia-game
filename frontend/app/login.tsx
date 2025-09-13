import React, { useState } from "react";
import { View, Text, TextInput, TouchableOpacity, Alert } from "react-native";
import { useRouter } from "expo-router";

export default function LoginScreen() {
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = () => {
    if (username && password) {
      Alert.alert("Login Success", `Welcome ${username}`);
      router.push("/game");
    } else {
      Alert.alert("Error", "Please enter username and password");
    }
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
