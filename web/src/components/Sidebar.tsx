import { ReactNode, ReactText } from "react"
import { IconType } from 'react-icons';
import {
  IoBugOutline,
  IoGitBranchOutline,
  IoHomeOutline,
  IoReceiptOutline,
  IoSettingsOutline,
  IoShieldCheckmarkOutline
} from "react-icons/io5";
import {
  Link
} from "react-router-dom";
import {
  Box,
  BoxProps,
  CloseButton,
  Flex,
  Icon,
  Text,
  useColorModeValue,
  useDisclosure
} from '@chakra-ui/react';
import { routes } from "../routes";

export default function Sidebar({ children }: { children: ReactNode }) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  return (
    <Box minH="100vh" bg={useColorModeValue('gray.100', 'gray.900')}>
      <SidebarContent
        onClose={() => onClose}
        display={{ base: 'none', md: 'block' }}
      />
      <Box ml={{ base: 0, md: 60 }} p="4">
        {children}
      </Box>
    </Box>
  );
}

interface SidebarProps extends BoxProps {
  onClose: () => void;
}

const SidebarContent = ({ onClose, ...args }: SidebarProps) => {
  return (
    <Box
      bg={useColorModeValue('white', 'gray.900')}
      borderRight="1px"
      borderRightColor={useColorModeValue('gray.200', 'gray.700')}
      w={{ base: 'full', md: 60 }}
      pos="fixed"
      h="full"
      {...args}>
      <Flex h="20" alignItems="center" mx="8" justifyContent="space-between">
        <Text color="git.primary" fontSize="2xl" fontFamily="monospace" fontWeight="bold">
          Watchdog
        </Text>
        <CloseButton display={{ base: 'flex', md: 'none' }} onClick={onClose} />
      </Flex>
      {routes.map((route) => (
        !route.hide && <NavItem key={route.title} route={route.path} icon={route.icon}>
          {route.title}
        </NavItem>
      ))}
    </Box>
  )
}

const NavItem = ({ children, route, icon, ...args }: { children: ReactText, route: string, icon?: IconType }) => {
  return (
    <Link to={route} style={{ textDecoration: 'none' }}>
      <Flex
        align="center"
        p="4"
        mx="4"
        borderRadius="lg"
        role="group"
        cursor="pointer"
        _hover={{
          bg: 'git.primary',
          color: 'white',
        }}
        {...args}>
          {icon && (
            <Icon
              mr="4"
              fontSize="16"
              _groupHover={{
                color: 'white',
              }}
              as={icon}
            />
          )}
          {children}
        </Flex>
    </Link>
  )
}
