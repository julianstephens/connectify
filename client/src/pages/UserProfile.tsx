import { Profile } from "@/components/profile";
import { useAuth0 } from "@auth0/auth0-react";
import { Button, Heading, Input, VStack } from "@chakra-ui/react";

const UserProfilePage = () => {
  const { user } = useAuth0();

  return (
    <>
      <Heading textAlign="center" mb="10">
        {user?.nickname ?? "User"}&#39;s Profile
      </Heading>
      <VStack gap="6" w="fit" mx="auto">
        <Profile size="2xl" />
        <VStack minW="300px" gap="3">
          <Input fontWeight="bold" value={user?.name} disabled={true} />
          <Input value={user?.email} disabled={true} />
        </VStack>
        <Button mt="10" colorPalette="red" variant="outline">
          Logout
        </Button>
      </VStack>
    </>
  );
};

export default UserProfilePage;
