import { useAuth0 } from "@auth0/auth0-react";
import { Box, Flex, Heading, Link, Menu, Portal } from "@chakra-ui/react";
import type { PropsWithChildren } from "react";
import { Profile } from "./profile";

export const PageLayout = ({ children }: PropsWithChildren) => {
  const { logout } = useAuth0();
  return (
    <Box padding="10">
      <Flex w="full" justify="space-between" align="center" mb="20">
        <Link href="/">
          <Heading size="2xl">Connectify</Heading>
        </Link>
        <Menu.Root positioning={{ placement: "bottom-end" }}>
          <Menu.Trigger rounded="full" focusRing="outside">
            <Profile />
          </Menu.Trigger>
          <Portal>
            <Menu.Positioner>
              <Menu.Content>
                <Menu.Item value="settings">
                  <Link href="/settings">Settings</Link>
                </Menu.Item>
                <Menu.Item value="logout" onClick={() => logout()}>
                  Logout
                </Menu.Item>
              </Menu.Content>
            </Menu.Positioner>
          </Portal>
        </Menu.Root>
      </Flex>
      {children}
    </Box>
  );
};
