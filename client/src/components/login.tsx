import { useAuth0 } from "@auth0/auth0-react";
import { Button, Flex, Heading } from "@chakra-ui/react";

export const Login = () => {
  const { isLoading: authLoading, loginWithRedirect } = useAuth0();
  return (
    <Flex direction="column" align="center" justify="center" gap={4}>
      <Heading>Welcome to Connectify</Heading>
      <Button onClick={() => loginWithRedirect()} disabled={authLoading}>
        {authLoading ? "Loading..." : "Login with Auth0"}
      </Button>
    </Flex>
  );
};
