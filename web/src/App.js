import "./App.css";

import { BrowserRouter as Router, Route } from "react-router-dom";
import { ChakraProvider } from "@chakra-ui/react";
import { Container } from "@chakra-ui/react";
import { Text } from "@chakra-ui/react";
import Nav from "./components/Nav";
import Stonks from "./components/Stonks";

function App() {
  return (
    <ChakraProvider>
      <Router>
        <Container className="app-header" centerContent>
          <Text fontSize="3xl" color="">
            Trade
          </Text>
          <Nav />
        </Container>
        <Container>
          <Route path="/">
            <Stonks />
          </Route>
        </Container>
      </Router>
    </ChakraProvider>
  );
}

export default App;
