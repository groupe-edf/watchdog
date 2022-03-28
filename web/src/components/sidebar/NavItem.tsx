import { Flex, Icon, Link } from "@chakra-ui/react"
import { ReactNode } from "react"
import { IconType } from "react-icons"
import { NavLink } from "react-router-dom"

interface NavItemProps {
  children: ReactNode,
  icon?: IconType,
  onClick?: () => void,
  route?: string,
  size: string
}

export const ItemContent = ({ children, icon, onClick, route, size, ...rest }: NavItemProps) => {
  return (
    <Flex
      _hover={{ background: 'brand.100', color: 'white'}}
      align="center"
      padding="2"
      width="100%"
      cursor="pointer"
      borderRadius="md"
      role="group"
      {...rest}>
      {icon && (
        <Icon marginRight="4" _groupHover={{ color: 'white' }} as={icon} />
      )}
      {children}
    </Flex>
  )
}

export const NavItem = ({ children, icon, onClick, route, size, ...rest }: NavItemProps) => {
  return (
    <>
    {route != null ? (
      <Link as={(props: any) => (
        <NavLink
          {...props}
          style={({ isActive }) => {
            return {
              fontWeight: isActive ? 700 : 400,
              textDecoration: "none"
            }
          }}
        />
      )}
        to={route}
        width="100%">
        <ItemContent children={children} icon={icon} size={size} {...rest} />
      </Link>
    ) : (
      <ItemContent children={children} icon={icon} size={size} onClick={onClick} {...rest}/>
    )}
    </>
  )
}
