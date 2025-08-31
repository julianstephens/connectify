import { Login } from "@/components/login";
import { useAuth0 } from "@auth0/auth0-react";
import { Flex } from "@chakra-ui/react";
import { Navigate } from "react-router-dom";

const LandingPage = () => {
  const { isAuthenticated } = useAuth0();

  return (
    <Flex id="landingPage" justify="center" align="center" height="full">
      {!isAuthenticated ? <Login /> : <Navigate to="/settings" />}
    </Flex>
  );
};

export default LandingPage;
