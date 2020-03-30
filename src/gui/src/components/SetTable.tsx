import React from 'react'
import { IRedisEntry, ILayoutModelState, connect } from 'umi';
import { Table, Button } from 'antd'
import { ColumnProps } from 'antd/es/table';
import { DeleteOutlined } from '@ant-design/icons';
import u from '@/utils/u';
import TableComponent from './TableComponent';

interface IPageProps {
    configs: any;
    entries: [];
}

class SetTable extends TableComponent<IPageProps> {
    _columns: ColumnProps<IRedisEntry>[] = [
        {
            title: 'Member',
            dataIndex: 'Value',
            key: 'Value',
            // onCell: this.onCell,
            className: "pointer",
            defaultSortOrder: "ascend",
            sorter: (a, b) => b.Value.localeCompare(a.Value),
            ...this.getColumnSearchProps('Value'),
        },
        {
            title: 'Action',
            dataIndex: '',
            key: 'x',
            width: 70,
            className: "ar",
            render: () => <Button type="danger" size="small" title="Delete"><DeleteOutlined /></Button>,
        },
    ];

    render() {
        const { configs, entries } = this.props;
        let pageSize = 15;
        if (!u.isNoW(configs) && !u.isNoW(configs.PageSize) && !u.isNoW(configs.PageSize.SubList)) {
            pageSize = configs.PageSize.SubList;
        }

        return (
            <Table<IRedisEntry>
                rowKey="Key"
                className="sublist"
                columns={this._columns}
                dataSource={entries}
                pagination={{ pageSize: pageSize, hideOnSinglePage: true }}
                bordered={true}
                size="small"
            />
        );
    }
}

export default connect(({ layout, }: { layout: ILayoutModelState; }) => ({
    configs: layout.Configs,
}))(SetTable);