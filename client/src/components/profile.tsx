import { useAuth0 } from "@auth0/auth0-react";
import { Avatar } from "@chakra-ui/react";

export const Profile = ({ size = "md" }: { size?: string }) => {
  const { user } = useAuth0();
  return (
    <Avatar.Root size={size as any}>
      <Avatar.Fallback name={user?.name} />
      <Avatar.Image src={user?.picture} />
    </Avatar.Root>
  );
};
