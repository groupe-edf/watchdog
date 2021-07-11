import { Button, Flex } from "@chakra-ui/react";
import { Component, ReactNode } from "react"

interface PaginationProps {
  currentPage: number,
  itemsPerPage: number,
  offset?: number,
  onChangePage: (props: any) => void,
  pagesToShow: number,
  totalItems: number,
}

class Pagination extends Component<PaginationProps, {
  currentPage: number,
  itemsPerPage: number,
  offset?: number,
  pagesToShow: number,
  totalItems: number
}> {
  constructor(props: PaginationProps) {
    super(props);
    this.state = {
      currentPage: 1,
      itemsPerPage: 10,
      offset: 0,
      pagesToShow: 5,
      totalItems: 0
    };
  }
  componentDidMount() {
    this.setState({
      currentPage: this.props.currentPage || 1,
      itemsPerPage: this.props.itemsPerPage || 10,
      offset: this.props.offset || 0,
      totalItems: this.props.totalItems || 0,
    });
  }
  componentWillReceiveProps(nextProps: PaginationProps) {
    this.setState({
      itemsPerPage: nextProps.itemsPerPage || 10,
      offset: nextProps.offset || 0,
      totalItems: nextProps.totalItems || 0,
      pagesToShow: nextProps.pagesToShow || 5
    });
  }
  componentDidUpdate(prevProps: PaginationProps, prevState: any) {
    if (
      this.state.itemsPerPage !== prevState.itemsPerPage ||
      this.state.totalItems !== prevState.totalItems
    ) {
      this.setPage(this.state.currentPage);
    }
  }
  setPage(page: number) {
    let { totalItems, itemsPerPage } = this.state;
    let totalPages = Math.ceil(totalItems / itemsPerPage)
    if (page < 1) {
      page = 1;
    } else if (page > totalPages) {
      page = totalPages;
    }
    this.setState({
      currentPage: page
    });
    var startIndex = page > 0 ? (page - 1) * itemsPerPage : 0
    var endIndex = Math.min(startIndex + itemsPerPage - 1, totalItems - 1)
    this.props.onChangePage({
      itemsPerPage,
      page,
      startIndex,
      totalPages,
      totalItems,
      endIndex
    });
  }
  getPager() {
    let { currentPage, pagesToShow, totalItems, itemsPerPage } = this.state;
    let totalPages = Math.ceil(totalItems / itemsPerPage)
    let pages = [];
    let startFromNumber;
    if (totalPages <= pagesToShow) {
      startFromNumber = 1;
      pagesToShow = totalPages;
    } else {
      if (currentPage <= Math.ceil(pagesToShow / 2)) {
        startFromNumber = 1;
      } else if (
        currentPage + Math.floor((pagesToShow - 1) / 2) >=
        totalPages
      ) {
        startFromNumber = totalPages - (pagesToShow - 1);
      } else {
        startFromNumber = currentPage - Math.floor(pagesToShow / 2);
      }
    }
    for (let i = 1; i <= pagesToShow; i++) {
      pages.push(startFromNumber++);
    }
    return {
      currentPage,
      totalPages,
      pages
    };
  }
  render() {
    let pager = this.getPager();
    const PageButton = ({ active, children, disabled, onClick, ...rest }: { active?: boolean, children: ReactNode, disabled?: boolean, onClick?: () => void }) => {
      const activeStyle = {
        background: "brand.100",
        color: "white",
      };
      return (
        <Button
          paddingX={4}
          paddingY={2}
          marginX={1}
          onClick={onClick}
          disabled={disabled}
          {...(active === true && activeStyle)}
          {...rest}>
          {children}
        </Button>
      )
    }
    return (
      <Flex
        alignItems="center"
        justifyContent="center"
        paddingTop={2}>
        <PageButton
          disabled={pager.currentPage === 1 ? true : false}
          onClick={() => this.setPage(1)}>
          Start
        </PageButton>
        <PageButton
          disabled={pager.currentPage === 1 ? true : false}
          onClick={() => this.setPage(pager.currentPage - 1)}>
          Previous
        </PageButton>
        {pager.pages.map((page, index) => (
          <PageButton
            key={index}
            active={pager.currentPage === page ? true : false}
            onClick={() => this.setPage(page)}>
            {page}
          </PageButton>
        ))}
        <PageButton
          disabled={pager.currentPage === pager.totalPages ? true : false}
          onClick={() => this.setPage(pager.currentPage + 1)}>
          Next
        </PageButton>
        <PageButton
          disabled={pager.currentPage === pager.totalPages ? true : false}
          onClick={() => this.setPage(pager.totalPages)}>
          End
        </PageButton>
      </Flex>
    )
  }
}

export { Pagination }
