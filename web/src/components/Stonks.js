import {
  Container,
  Flex,
  ListItem,
  Text,
  UnorderedList,
} from "@chakra-ui/react";
import { useState } from "react";

const Stonks = () => {
  const [symbols, setSymbols] = useState([
    { name: "JOBR", price: 2000 },
    { name: "GME", price: 30000 },
  ]);

  return (
    <Flex>
      <UnorderedList>
        {symbols.map((symbol) => (
          <ListItem key={symbol.name}>{symbol.name}</ListItem>
        ))}
      </UnorderedList>
      <Container>hello stonk</Container>
    </Flex>
  );
};

export default Stonks;
