import { Breadcrumb, BreadcrumbItem, BreadcrumbLink } from "@chakra-ui/react";
import { Link } from "react-router-dom";
const Nav = () => (
  <Breadcrumb separator="-">
    <BreadcrumbItem>
      <BreadcrumbLink as={Link} to="/">
        Home
      </BreadcrumbLink>
    </BreadcrumbItem>

    <BreadcrumbItem>
      <BreadcrumbLink as={Link} to="/symbols">
        Symbols
      </BreadcrumbLink>
    </BreadcrumbItem>

    <BreadcrumbItem isCurrentPage>
      <BreadcrumbLink as={Link} to="/portfolio">
        Portfolio
      </BreadcrumbLink>
    </BreadcrumbItem>
  </Breadcrumb>
);

export default Nav;
