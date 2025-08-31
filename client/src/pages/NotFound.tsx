import { Button, Flex, Heading, Link } from "@chakra-ui/react";

const NotFoundPage = () => (
  <Flex
    direction="column"
    gap="4"
    justify="center"
    align="center"
    height="full"
  >
    <Heading size="2xl">404 - Not Found</Heading>
    <Link href="/">
      <Button>Go Home</Button>
    </Link>
  </Flex>
);

export default NotFoundPage;
