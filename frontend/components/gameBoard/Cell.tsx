import { TouchableOpacity, Text } from "react-native";
import React from "react";

type CellProps = {
  name: string;
  marked: boolean;
  onPress: () => void;
};

const Cell = ({ name, marked, onPress }: CellProps) => (
  <TouchableOpacity
    onPress={onPress}
    style={{
      width: 64,
      height: 64,
      borderWidth: 1,
      alignItems: "center",
      justifyContent: "center",
      backgroundColor: marked ? "green" : "white",
      padding: 2,
    }}
  >
    <Text>{name}</Text>
  </TouchableOpacity>
);

export default Cell;
