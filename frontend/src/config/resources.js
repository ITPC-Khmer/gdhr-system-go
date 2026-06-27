// Resource definitions that drive the generic CRUD list & form components.
//
// Each resource declares:
//   endpoint          API path (relative to /api)
//   title / singular  display labels
//   icon              sidebar/icon name
//   pk                primary-key field name
//   pkEditable        whether the PK can be typed (true) — disabled when editing
//   searchPlaceholder placeholder for the free-text "search any" box
//   columns[]         table columns: { key, label, type?, mono?, width? }
//   filters[]         filter controls: { key, label, type, options? }
//   fields[]          form fields: { key, label, type, required?, span?, options? }
//
// type values: 'text' | 'number' | 'bool' | 'textarea' | 'select'

const ACTIVE_FILTER = {
  key: 'active',
  label: 'Status',
  type: 'select',
  options: [
    { value: '', label: 'All' },
    { value: '1', label: 'Active' },
    { value: '0', label: 'Inactive' },
  ],
}

export const resources = {
  institutes: {
    endpoint: '/institutes',
    title: 'Institutes',
    singular: 'Institute',
    icon: 'building',
    pk: 'id',
    pkEditable: true,
    searchPlaceholder: 'Search name, code, id…',
    columns: [
      {
        key: 'profile',
        type: 'profile',
        label: 'Institute',
        avatarIcon: 'building',
        title: 'name',
        subtitle: 'name_short',
        tags: [
          { key: 'source_table', class: 'bg-slate-100 text-slate-600' },
        ],
      },
      { key: 'parent_name', label: 'Parent', truncate: true },
      { key: 'level_type', label: 'Level' },
      { key: 'active', label: 'Active', type: 'bool' },
    ],
    filters: [
      { key: 'level_type', label: 'Level type', type: 'number' },
      { key: 'source_table', label: 'Source table', type: 'text' },
      { key: 'parent_id', label: 'Parent ID', type: 'text' },
      ACTIVE_FILTER,
    ],
    fields: [
      { key: 'id', label: 'ID', type: 'text', required: true },
      { key: 'old_id', label: 'Old ID', type: 'text' },
      { key: 'parent_id', label: 'Parent ID', type: 'text' },
      { key: 'level_type', label: 'Level type', type: 'number' },
      { key: 'old_parent_id', label: 'Old parent ID', type: 'number' },
      { key: 'name', label: 'Name', type: 'text', required: true, span: 2 },
      { key: 'name_short', label: 'Short name', type: 'text', span: 2 },
      { key: 'source_table', label: 'Source table', type: 'text' },
      { key: 'source_table_id', label: 'Source table ID', type: 'number' },
      { key: 'active', label: 'Active', type: 'bool' },
    ],
  },

  staffs: {
    endpoint: '/staffs',
    title: 'Staffs',
    singular: 'Staff',
    icon: 'users',
    pk: 'uid',
    pkEditable: true,
    searchPlaceholder: 'Search name, staff no., phone, email…',
    baseFilter: { status_id: 2 }, // list shows only active staff (status_id = 2)
    columns: [
      {
        key: 'profile',
        type: 'profile',
        label: 'Staff',
        image: 'photo_path',
        title: ['rank_name_short', 'surname_kh', 'name_kh'],
        subtitle: ['surname_en', 'name_en'],
        badge: 'staff_number',
        tags: [
          { key: 'position_name_short', class: 'bg-indigo-100 text-indigo-700' },
        ],
      },
      { key: 'institute_hierarchy', label: 'Institute', type: 'chain', separator: ' → ', dropLast: 1, reverse: true },
      { key: 'department_name_short', label: 'Department' },
    ],
    filters: [
      { key: 'rank_id', label: 'Rank ID', type: 'number' },
      { key: 'position_id', label: 'Position ID', type: 'number' },
      { key: 'department_id', label: 'Department ID', type: 'number' },
      { key: 'institute_id', label: 'Institute ID', type: 'text' },
      {
        key: 'gender',
        label: 'Gender',
        type: 'select',
        options: [
          { value: '', label: 'All' },
          { value: '1', label: 'Male' },
          { value: '2', label: 'Female' },
        ],
      },
    ],
    fields: [
      { key: 'uid', label: 'UID', type: 'text', required: true, span: 2 },
      { key: 'staff_number', label: 'Staff number', type: 'text' },
      { key: 'staff_type_id', label: 'Staff type ID', type: 'number' },
      { key: 'surname_kh', label: 'Surname (KH)', type: 'text' },
      { key: 'name_kh', label: 'Name (KH)', type: 'text' },
      { key: 'surname_en', label: 'Surname (EN)', type: 'text' },
      { key: 'name_en', label: 'Name (EN)', type: 'text' },
      {
        key: 'gender',
        label: 'Gender',
        type: 'select',
        options: [
          { value: 0, label: '—' },
          { value: 1, label: 'Male' },
          { value: 2, label: 'Female' },
        ],
      },
      { key: 'phone', label: 'Phone', type: 'text' },
      { key: 'email', label: 'Email', type: 'text' },
      { key: 'place_of_birth', label: 'Place of birth', type: 'text' },
      { key: 'nationality', label: 'Nationality', type: 'text' },
      { key: 'cityzen_card_number', label: 'Citizen card no.', type: 'text' },
      { key: 'address', label: 'Address', type: 'textarea', span: 2 },
      { key: 'photo_path', label: 'Photo', type: 'image', span: 2 },
      { key: 'rank_id', label: 'Rank ID', type: 'number' },
      { key: 'rank_name_short', label: 'Rank name (short)', type: 'text' },
      { key: 'position_id', label: 'Position ID', type: 'number' },
      { key: 'position_name_short', label: 'Position name (short)', type: 'text' },
      { key: 'other_position_id', label: 'Other position ID', type: 'number' },
      { key: 'general_commissariat_id', label: 'Gen. commissariat ID', type: 'number' },
      { key: 'general_commissariat_name_short', label: 'Gen. commissariat (short)', type: 'text' },
      { key: 'department_id', label: 'Department ID', type: 'number' },
      { key: 'department_name_short', label: 'Department (short)', type: 'text' },
      { key: 'office_id', label: 'Office ID', type: 'number' },
      { key: 'office_name_short', label: 'Office (short)', type: 'text' },
      { key: 'sector_id', label: 'Sector ID', type: 'number' },
      { key: 'sector_name_short', label: 'Sector (short)', type: 'text' },
      { key: 'institute_id', label: 'Institute ID', type: 'text', span: 2 },
      { key: 'status_id', label: 'Status ID', type: 'number' },
      { key: 'status_name', label: 'Status name', type: 'text' },
    ],
  },

  ranks: {
    endpoint: '/ranks',
    title: 'Ranks',
    singular: 'Rank',
    icon: 'shield',
    pk: 'rank_id',
    pkEditable: true,
    searchPlaceholder: 'Search rank name…',
    columns: [
      { key: 'rank_id', label: 'ID' },
      { key: 'rank_name', label: 'Name (KH)' },
      { key: 'rank_name_short', label: 'Short' },
      { key: 'rank_name_en', label: 'Name (EN)' },
      { key: 'rank_order', label: 'Order' },
      { key: 'promote_period', label: 'Promote (yr)' },
      { key: 'active', label: 'Active', type: 'bool' },
    ],
    filters: [
      { key: 'position_base_id', label: 'Position base ID', type: 'number' },
      { key: 'rank_order', label: 'Rank order', type: 'number' },
      ACTIVE_FILTER,
    ],
    fields: [
      { key: 'rank_id', label: 'Rank ID', type: 'number', required: true },
      { key: 'rank_order', label: 'Rank order', type: 'number' },
      { key: 'rank_name', label: 'Rank name (KH)', type: 'text', span: 2 },
      { key: 'rank_name_short', label: 'Short name (KH)', type: 'text' },
      { key: 'rank_name_en', label: 'Rank name (EN)', type: 'text' },
      { key: 'rank_name_short_en', label: 'Short name (EN)', type: 'text' },
      { key: 'position_base_id', label: 'Position base ID', type: 'number' },
      { key: 'promote_period', label: 'Promote period (years)', type: 'number' },
      { key: 'active', label: 'Active', type: 'bool' },
    ],
  },

  positions: {
    endpoint: '/positions',
    title: 'Positions',
    singular: 'Position',
    icon: 'briefcase',
    pk: 'position_id',
    pkEditable: true,
    searchPlaceholder: 'Search position name…',
    columns: [
      { key: 'position_id', label: 'ID' },
      { key: 'position_name', label: 'Name (KH)' },
      { key: 'position_name_short', label: 'Short' },
      { key: 'position_base_id', label: 'Position base' },
      { key: 'rank_base_id', label: 'Rank base' },
    ],
    filters: [
      { key: 'position_base_id', label: 'Position base ID', type: 'number' },
      { key: 'rank_base_id', label: 'Rank base ID', type: 'number' },
    ],
    fields: [
      { key: 'position_id', label: 'Position ID', type: 'number', required: true },
      { key: 'position_name', label: 'Position name (KH)', type: 'text', span: 2 },
      { key: 'position_name_short', label: 'Short name', type: 'text' },
      { key: 'position_base_id', label: 'Position base ID', type: 'number' },
      { key: 'rank_base_id', label: 'Rank base ID', type: 'number' },
      { key: 'description', label: 'Description', type: 'textarea', span: 2 },
    ],
  },

  holidays: {
    endpoint: '/holidays',
    title: 'Holidays',
    singular: 'Holiday',
    icon: 'calendar',
    pk: 'holiday_id',
    pkEditable: false, // server generates the UUID
    searchPlaceholder: 'Search holiday name…',
    columns: [
      { key: 'name', label: 'Name' },
      { key: 'date', label: 'Date', type: 'date' },
      { key: 'from_date', label: 'From', type: 'date' },
      { key: 'to_date', label: 'To', type: 'date' },
      { key: 'is_recurring', label: 'Recurring', type: 'bool' },
      { key: 'holiday_id', label: 'ID', mono: true, truncate: true },
    ],
    filters: [
      {
        key: 'is_recurring',
        label: 'Recurring',
        type: 'select',
        options: [
          { value: '', label: 'All' },
          { value: '1', label: 'Recurring' },
          { value: '0', label: 'One-off' },
        ],
      },
    ],
    fields: [
      { key: 'name', label: 'Name', type: 'text', required: true, span: 2 },
      { key: 'date', label: 'Date', type: 'date', required: true },
      { key: 'is_recurring', label: 'Recurring', type: 'bool' },
      { key: 'from_date', label: 'From date', type: 'date' },
      { key: 'to_date', label: 'To date', type: 'date' },
    ],
  },

  leave_types: {
    endpoint: '/leave-types',
    title: 'Leave Types',
    singular: 'Leave Type',
    icon: 'list',
    pk: 'id',
    pkEditable: false,
    searchPlaceholder: 'Search key, name, description…',
    columns: [
      { key: 'id', label: 'ID' },
      { key: 'leave_key', label: 'Key', mono: true },
      { key: 'type_name', label: 'Name (KH)' },
      { key: 'type_name_s', label: 'Short' },
      { key: 'max_days', label: 'Max days' },
      { key: 'is_reset', label: 'Resets', type: 'bool' },
    ],
    filters: [
      {
        key: 'is_reset',
        label: 'Resets yearly',
        type: 'select',
        options: [
          { value: '', label: 'All' },
          { value: '1', label: 'Resets' },
          { value: '0', label: 'No reset' },
        ],
      },
    ],
    fields: [
      { key: 'leave_key', label: 'Leave key', type: 'text' },
      { key: 'type_name', label: 'Type name (KH)', type: 'text', required: true, span: 2 },
      { key: 'type_name_s', label: 'Short name', type: 'text', required: true },
      { key: 'max_days', label: 'Max days / year', type: 'number' },
      { key: 'is_reset', label: 'Resets each year', type: 'bool' },
      { key: 'description', label: 'Description', type: 'textarea', span: 2 },
    ],
  },

  leave_roles: {
    endpoint: '/leave-roles',
    title: 'Leave Roles',
    singular: 'Leave Role',
    icon: 'shield',
    pk: 'id',
    pkEditable: false,
    searchPlaceholder: 'Search name, leave type, staff type…',
    columns: [
      { key: 'id', label: 'ID' },
      { key: 'leave_type', label: 'Leave type' },
      { key: 'name', label: 'Name' },
      { key: 'min_duration', label: 'Min days' },
      { key: 'max_duration', label: 'Max days' },
      { key: 'approve_level', label: 'Approve level' },
      { key: 'staff_type', label: 'Staff type' },
    ],
    filters: [
      { key: 'leave_type_id', label: 'Leave type ID', type: 'number' },
      { key: 'leave_type', label: 'Leave type key', type: 'text' },
      { key: 'staff_type', label: 'Staff type', type: 'text' },
      { key: 'approve_level', label: 'Approve level', type: 'number' },
    ],
    fields: [
      { key: 'leave_type_id', label: 'Leave type ID', type: 'number' },
      { key: 'leave_type', label: 'Leave type key', type: 'text', required: true },
      { key: 'name', label: 'Name', type: 'text', required: true, span: 2 },
      { key: 'min_duration', label: 'Min duration (days)', type: 'number', required: true },
      { key: 'max_duration', label: 'Max duration (days)', type: 'number', required: true },
      { key: 'limit_days', label: 'Limit days', type: 'number', required: true },
      { key: 'min_duration_show', label: 'Min duration (show)', type: 'number', required: true },
      { key: 'approve_level', label: 'Approve level', type: 'number', required: true },
      { key: 'staff_type', label: 'Staff type', type: 'text', required: true },
    ],
  },

  leaves: {
    endpoint: '/leaves',
    title: 'Leave Requests',
    singular: 'Leave Request',
    icon: 'calendar',
    pk: 'id',
    pkEditable: false,
    detail: true, // enables the "view" link -> /leaves/:id (approval timeline)
    readOnly: true, // no add/edit/delete in-app; requests are filed via API, state changes via the approval workflow
    searchPlaceholder: 'Search staff, ref no., reason…',
    columns: [
      { key: 'id', label: 'ID' },
      { key: 'staff_id', label: 'Staff', mono: true },
      { key: 'leave_type_name', label: 'Type' },
      { key: 'start_date', label: 'Start', type: 'date' },
      { key: 'end_date', label: 'End', type: 'date' },
      { key: 'total_day', label: 'Days' },
      { key: 'status', label: 'Status' },
    ],
    filters: [
      { key: 'staff_id', label: 'Staff ID', type: 'text' },
      { key: 'leave_type_id', label: 'Leave type ID', type: 'number' },
      {
        key: 'status',
        label: 'Status',
        type: 'select',
        options: [
          { value: '', label: 'All' },
          { value: 'pending', label: 'Pending' },
          { value: 'approved', label: 'Approved' },
          { value: 'rejected', label: 'Rejected' },
        ],
      },
    ],
    // total_day, approval workflow and leave_year balance are handled server-side.
    fields: [
      { key: 'staff_id', label: 'Staff ID', type: 'text', required: true },
      { key: 'leave_type_id', label: 'Leave type ID', type: 'number', required: true },
      { key: 'ref_number', label: 'Reference number', type: 'text' },
      { key: 'start_date', label: 'Start date', type: 'date', required: true },
      { key: 'end_date', label: 'End date', type: 'date', required: true },
      { key: 'total_day', label: 'Total days (auto if 0)', type: 'number' },
      { key: 'phone', label: 'Phone', type: 'text' },
      {
        key: 'status',
        label: 'Status',
        type: 'select',
        options: [
          { value: 'pending', label: 'Pending' },
          { value: 'approved', label: 'Approved' },
          { value: 'rejected', label: 'Rejected' },
        ],
      },
      { key: 'reason', label: 'Reason', type: 'textarea', span: 2 },
      { key: 'reject_reason', label: 'Reject reason', type: 'textarea', span: 2 },
      { key: 'attachment', label: 'Attachment', type: 'textarea', span: 2 },
      { key: 'file_name', label: 'File name', type: 'text', span: 2 },
    ],
  },

  leave_approvals: {
    endpoint: '/leave-approvals',
    title: 'Leave Approvals',
    singular: 'Leave Approval',
    icon: 'shield',
    pk: 'id',
    pkEditable: false,
    searchPlaceholder: 'Search staff, role, institute…',
    // Per-approver task actions. method/path drive the API call ({id} -> row pk);
    // showWhen gates visibility by a row field; prompt asks for a body field.
    actions: [
      {
        key: 'approve',
        label: 'Approve',
        icon: 'check',
        variant: 'success',
        method: 'post',
        path: '/leave-approvals/{id}/approve',
        showWhen: { status: 'pending' },
        authApprover: true, // only the assigned approver (not admins) may act
        confirm: 'Approve this task? This advances or completes the request.',
      },
      {
        key: 'reject',
        label: 'Reject',
        icon: 'x',
        variant: 'danger',
        method: 'post',
        path: '/leave-approvals/{id}/reject',
        showWhen: { status: 'pending' },
        authApprover: true,
        prompt: { field: 'reject_reason', label: 'Reason for rejection (optional)' },
        confirm: 'Reject this task? This rejects the whole leave request.',
      },
    ],
    columns: [
      { key: 'id', label: 'ID' },
      { key: 'leave_id', label: 'Leave ID' },
      { key: 'staff_id', label: 'Approver', mono: true },
      { key: 'approve_level', label: 'Level' },
      { key: 'max_level', label: 'Max level' },
      { key: 'status', label: 'Status' },
    ],
    filters: [
      { key: 'leave_id', label: 'Leave ID', type: 'number' },
      { key: 'staff_id', label: 'Staff ID', type: 'text' },
      {
        key: 'l_type',
        label: 'Type',
        type: 'select',
        options: [
          { value: '', label: 'All' },
          { value: 'leave', label: 'Leave' },
          { value: 'mission', label: 'Mission' },
        ],
      },
      {
        key: 'status',
        label: 'Status',
        type: 'select',
        options: [
          { value: '', label: 'All' },
          { value: 'pending', label: 'Pending' },
          { value: 'approved', label: 'Approved' },
          { value: 'rejected', label: 'Rejected' },
        ],
      },
    ],
    fields: [
      { key: 'leave_id', label: 'Leave ID', type: 'number' },
      { key: 'staff_id', label: 'Approver staff ID', type: 'text' },
      {
        key: 'l_type',
        label: 'Type',
        type: 'select',
        options: [
          { value: 'leave', label: 'Leave' },
          { value: 'mission', label: 'Mission' },
        ],
      },
      { key: 'l_level', label: 'Level', type: 'number' },
      { key: 'approve_level', label: 'Approve level', type: 'number' },
      { key: 'max_level', label: 'Max level', type: 'number' },
      {
        key: 'status',
        label: 'Status',
        type: 'select',
        options: [
          { value: 'pending', label: 'Pending' },
          { value: 'approved', label: 'Approved' },
          { value: 'rejected', label: 'Rejected' },
        ],
      },
      { key: 'role_name', label: 'Role name', type: 'text' },
      { key: 'institute_id', label: 'Institute ID', type: 'text' },
      { key: 'show_type', label: 'Show type', type: 'number' },
      { key: 'moderator_type', label: 'Moderator type', type: 'number' },
    ],
  },

  leave_files: {
    endpoint: '/leave-files',
    title: 'Leave Files',
    singular: 'Leave File',
    icon: 'list',
    pk: 'id',
    pkEditable: false,
    searchPlaceholder: 'Search comment, staff…',
    columns: [
      { key: 'id', label: 'ID' },
      { key: 'leave_id', label: 'Leave ID' },
      { key: 'task_id', label: 'Task ID' },
      { key: 'message_type', label: 'Type' },
      { key: 'staff_id', label: 'Staff', mono: true },
      { key: 'status', label: 'Status' },
    ],
    filters: [
      { key: 'leave_id', label: 'Leave ID', type: 'number' },
      { key: 'task_id', label: 'Task ID', type: 'number' },
      {
        key: 'message_type',
        label: 'Message type',
        type: 'select',
        options: [
          { value: '', label: 'All' },
          { value: 'description', label: 'Description' },
          { value: 'chat', label: 'Chat' },
        ],
      },
    ],
    fields: [
      { key: 'leave_id', label: 'Leave ID', type: 'number', required: true },
      { key: 'task_id', label: 'Task ID', type: 'number', required: true },
      {
        key: 'message_type',
        label: 'Message type',
        type: 'select',
        options: [
          { value: 'description', label: 'Description' },
          { value: 'chat', label: 'Chat' },
        ],
      },
      { key: 'parent_id', label: 'Parent ID', type: 'number' },
      { key: 'staff_id', label: 'Staff ID', type: 'text' },
      { key: 'added_by', label: 'Added by', type: 'text' },
      { key: 'last_updated_by', label: 'Last updated by', type: 'text' },
      {
        key: 'status',
        label: 'Status',
        type: 'select',
        options: [
          { value: 'sent', label: 'Sent' },
          { value: 'read', label: 'Read' },
        ],
      },
      { key: 'comment', label: 'Comment', type: 'textarea', span: 2 },
    ],
  },

  leave_years: {
    endpoint: '/leave-years',
    title: 'Leave Balances',
    singular: 'Leave Balance',
    icon: 'calendar',
    pk: 'id',
    pkEditable: false,
    searchPlaceholder: 'Search staff…',
    columns: [
      { key: 'id', label: 'ID' },
      { key: 'staff_id', label: 'Staff', mono: true },
      { key: 'leave_type_id', label: 'Type ID' },
      { key: 'l_year', label: 'Year' },
      { key: 'total_day', label: 'Used' },
      { key: 'max_days', label: 'Max' },
      { key: 'leave_remaining', label: 'Remaining' },
    ],
    filters: [
      { key: 'staff_id', label: 'Staff ID', type: 'text' },
      { key: 'leave_type_id', label: 'Leave type ID', type: 'number' },
      { key: 'l_year', label: 'Year', type: 'number' },
    ],
    fields: [
      { key: 'staff_id', label: 'Staff ID', type: 'text', required: true },
      { key: 'leave_type_id', label: 'Leave type ID', type: 'number', required: true },
      { key: 'l_year', label: 'Year', type: 'number', required: true },
      { key: 'total_day', label: 'Days used', type: 'number' },
      { key: 'max_days', label: 'Max days', type: 'number' },
      { key: 'leave_remaining', label: 'Remaining', type: 'number' },
      { key: 'is_reset', label: 'Reset', type: 'bool' },
    ],
  },

  staff_institute_roles: {
    endpoint: '/staff-institute-roles',
    title: 'Manage Roles',
    singular: 'Manage Role',
    icon: 'shield',
    pk: 'id',
    pkEditable: false,
    searchPlaceholder: 'Search staff, institute, role…',
    columns: [
      { key: 'id', label: 'ID' },
      { key: 'staff_id', label: 'Staff', mono: true },
      { key: 'institute_id', label: 'Institute ID', mono: true },
      { key: 'institute_name', label: 'Institute', truncate: true },
      { key: 'role', label: 'Role' },
    ],
    filters: [
      { key: 'staff_id', label: 'Staff ID', type: 'text' },
      { key: 'institute_id', label: 'Institute ID', type: 'text' },
      {
        key: 'role',
        label: 'Role',
        type: 'select',
        options: [
          { value: '', label: 'All' },
          { value: 'admin', label: 'Admin' },
          { value: 'moderator', label: 'Moderator' },
          { value: 'approval', label: 'Approval' },
        ],
      },
    ],
    fields: [
      { key: 'staff_id', label: 'Staff ID', type: 'text', required: true },
      { key: 'institute_id', label: 'Institute ID', type: 'text', required: true },
      {
        key: 'role',
        label: 'Role',
        type: 'select',
        required: true,
        options: [
          { value: 'admin', label: 'Admin' },
          { value: 'moderator', label: 'Moderator' },
          { value: 'approval', label: 'Approval' },
        ],
      },
    ],
  },
}

export function getResource(key) {
  return resources[key]
}

export const resourceList = Object.entries(resources).map(([key, r]) => ({
  key,
  title: r.title,
  icon: r.icon,
}))
